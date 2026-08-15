package application

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/chengkz2023/My-GVA/server/internal/modules/business/file/domain"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
)

// MaxUploadBytes 单个文件上传大小上限（10 MB）。
const MaxUploadBytes = 10 << 20

// blockedUploadExts 拒绝可被同源直接执行/渲染的危险扩展名，防止存储型 XSS。
var blockedUploadExts = map[string]bool{
	".html": true, ".htm": true, ".svg": true, ".js": true, ".mjs": true,
	".php": true, ".jsp": true, ".asp": true, ".aspx": true,
}

func validateUploadType(filename string, data []byte) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if blockedUploadExts[ext] {
		return apperrors.WithMessage(apperrors.Validation, "file type not allowed")
	}
	if ct := http.DetectContentType(data); strings.HasPrefix(ct, "text/html") {
		return apperrors.WithMessage(apperrors.Validation, "file type not allowed")
	}
	return nil
}

type Service struct {
	repo      domain.Repository
	storePath string
	accessURL string
}

func NewService(repo domain.Repository, storePath string) *Service {
	return &Service{repo: repo, storePath: storePath, accessURL: storePath}
}

func (s *Service) Upload(ctx context.Context, header FileHeader, classID int) (UploadResponse, error) {
	if s.repo == nil {
		return UploadResponse{}, apperrors.WithMessage(apperrors.Internal, "file repository unavailable")
	}
	if len(header.Data) == 0 {
		return UploadResponse{}, apperrors.WithMessage(apperrors.Validation, "file is empty")
	}
	if len(header.Data) > MaxUploadBytes {
		return UploadResponse{}, apperrors.WithMessage(apperrors.Validation, "file too large")
	}
	if err := validateUploadType(header.Filename, header.Data); err != nil {
		return UploadResponse{}, err
	}

	filename := s.safeFilename(header.Filename)
	storePath := filepath.Join(s.storePath, filename)

	if err := writeFile(storePath, header.Data); err != nil {
		return UploadResponse{}, apperrors.New(apperrors.Internal, 0, "save file failed", err)
	}

	url := s.accessURL + "/" + filename
	tag := strings.TrimPrefix(filepath.Ext(header.Filename), ".")

	record, err := s.repo.Create(ctx, domain.File{
		Name:    header.Filename,
		ClassID: classID,
		URL:     url,
		Tag:     tag,
		Key:     filename,
	})
	if err != nil {
		// 落盘成功但入库失败：回滚已写入的文件，避免产生孤儿文件
		_ = removeFile(storePath)
		return UploadResponse{}, apperrors.New(apperrors.Internal, 0, "save file record failed", err)
	}
	return UploadResponse{File: fromDomain(record)}, nil
}

func (s *Service) List(ctx context.Context, query domain.ListQuery) (ListResponse, error) {
	if s.repo == nil {
		return ListResponse{List: []FileResponse{}}, nil
	}
	result, err := s.repo.List(ctx, query)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return ListResponse{List: []FileResponse{}}, nil
	}
	if err != nil {
		return ListResponse{}, apperrors.New(apperrors.Internal, 0, "list files failed", err)
	}
	items := make([]FileResponse, 0, len(result.List))
	for _, f := range result.List {
		items = append(items, fromDomain(f))
	}
	return ListResponse{
		List: items, Total: result.Total, Page: result.Page, PageSize: result.PageSize,
	}, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "file repository unavailable")
	}
	file, err := s.repo.Delete(ctx, id)
	if err != nil {
		return apperrors.New(apperrors.Internal, 0, "delete file record failed", err)
	}
	if file.Key == "" {
		return nil
	}
	// 校验 key 为纯文件名，拒绝含路径分隔符或 .. 的脏数据，避免越界删除
	if filepath.Base(file.Key) != file.Key || strings.Contains(file.Key, "..") {
		return nil
	}
	if err := removeFile(filepath.Join(s.storePath, file.Key)); err != nil {
		// 文件本就不存在视为删除成功；真正失败（权限等）才返回错误
		if !os.IsNotExist(err) {
			return apperrors.New(apperrors.Internal, 0, "delete file failed", err)
		}
	}
	return nil
}

func (s *Service) Update(ctx context.Context, id uint, name string, tag string) (FileResponse, error) {
	if s.repo == nil {
		return FileResponse{}, apperrors.WithMessage(apperrors.Internal, "file repository unavailable")
	}
	file, err := s.repo.Update(ctx, id, name, tag)
	if err != nil {
		return FileResponse{}, apperrors.New(apperrors.Internal, 0, "update file failed", err)
	}
	return fromDomain(file), nil
}

func (s *Service) safeFilename(original string) string {
	return sanitizeFilename(original)
}

func fromDomain(f domain.File) FileResponse {
	return FileResponse{
		ID:      f.ID,
		Name:    f.Name,
		ClassID: f.ClassID,
		URL:     f.URL,
		Tag:     f.Tag,
		Key:     f.Key,
	}
}

type FileHeader struct {
	Filename string
	Data     []byte
}

func mapFiles(files []domain.File) []FileResponse {
	items := make([]FileResponse, 0, len(files))
	for _, f := range files {
		items = append(items, fromDomain(f))
	}
	return items
}

// These functions are in a separate file to avoid duplicating utils
func writeFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := ensureDir(dir); err != nil {
		return err
	}
	return writeBytesToFile(path, data)
}

func removeFile(path string) error {
	return remove(path)
}

func ensureDir(dir string) error {
	return mkdirAll(dir)
}
