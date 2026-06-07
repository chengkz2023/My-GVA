package application

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/domain"
	apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
)

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
	if file.Key != "" {
		_ = removeFile(filepath.Join(s.storePath, file.Key))
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
