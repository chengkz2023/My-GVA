package scaffold

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const defaultModuleRoot = "internal/modules"

var segmentPattern = regexp.MustCompile(`^[a-z][a-z0-9]*$`)
var routePattern = regexp.MustCompile(`^(/[a-z][a-z0-9-]*)+$`)

type ModuleOptions struct {
	RootDir string
	Group   string
	Name    string
	Route   string
	Force   bool
}

type ModuleResult struct {
	Dir   string
	Files []string
}

type moduleTemplateData struct {
	PackageName string
	ModulePath  string
	RoutePath   string
}

func GenerateLightModule(opts ModuleOptions) (ModuleResult, error) {
	opts = normalizeModuleOptions(opts)
	if err := validateModuleOptions(opts); err != nil {
		return ModuleResult{}, err
	}

	dir := filepath.Join(opts.RootDir, defaultModuleRoot, opts.Group, opts.Name)
	if exists, err := pathExists(dir); err != nil {
		return ModuleResult{}, err
	} else if exists && !opts.Force {
		return ModuleResult{}, fmt.Errorf("%s already exists", dir)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ModuleResult{}, err
	}

	data := moduleTemplateData{
		PackageName: opts.Name,
		ModulePath:  opts.Group + "/" + opts.Name,
		RoutePath:   opts.Route,
	}

	result := ModuleResult{Dir: dir}
	for _, file := range lightModuleFiles {
		path := filepath.Join(dir, file.name)
		if exists, err := pathExists(path); err != nil {
			return ModuleResult{}, err
		} else if exists && !opts.Force {
			return ModuleResult{}, fmt.Errorf("%s already exists", path)
		}

		content, err := renderGoTemplate(file.name, file.body, data)
		if err != nil {
			return ModuleResult{}, err
		}
		if err := os.WriteFile(path, content, 0o644); err != nil {
			return ModuleResult{}, err
		}
		result.Files = append(result.Files, path)
	}

	return result, nil
}

func normalizeModuleOptions(opts ModuleOptions) ModuleOptions {
	if opts.RootDir == "" {
		opts.RootDir = "."
	}
	opts.Group = strings.ToLower(strings.TrimSpace(opts.Group))
	if opts.Group == "" {
		opts.Group = "business"
	}
	opts.Name = strings.ToLower(strings.TrimSpace(opts.Name))
	opts.Route = strings.TrimSpace(opts.Route)
	if opts.Route == "" {
		opts.Route = "/" + opts.Name
	}
	if !strings.HasPrefix(opts.Route, "/") {
		opts.Route = "/" + opts.Route
	}
	return opts
}

func validateModuleOptions(opts ModuleOptions) error {
	if opts.Name == "" {
		return errors.New("module name is required")
	}
	if !segmentPattern.MatchString(opts.Group) {
		return fmt.Errorf("invalid group %q: use lowercase letters and digits, starting with a letter", opts.Group)
	}
	if !segmentPattern.MatchString(opts.Name) {
		return fmt.Errorf("invalid name %q: use lowercase letters and digits, starting with a letter", opts.Name)
	}
	if opts.Route == "/" {
		return errors.New("route must not be root")
	}
	if !routePattern.MatchString(opts.Route) {
		return fmt.Errorf("invalid route %q: use lowercase /segments with letters, digits, and hyphens", opts.Route)
	}
	return nil
}

func renderGoTemplate(name string, body string, data moduleTemplateData) ([]byte, error) {
	tpl, err := template.New(name).Parse(body)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format %s: %w", name, err)
	}
	return formatted, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

type lightModuleFile struct {
	name string
	body string
}

var lightModuleFiles = []lightModuleFile{
	{name: "dto.go", body: dtoTemplate},
	{name: "handler.go", body: handlerTemplate},
	{name: "handler_test.go", body: handlerTestTemplate},
	{name: "model.go", body: modelTemplate},
	{name: "module.go", body: moduleTemplate},
	{name: "repository.go", body: repositoryTemplate},
	{name: "service.go", body: serviceTemplate},
}

const dtoTemplate = `package {{.PackageName}}

type InfoResponse struct {
	Name string ` + "`json:\"name\"`" + `
	Message string ` + "`json:\"message\"`" + `
}
`

const handlerTemplate = `package {{.PackageName}}

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(group *gin.RouterGroup) {
	moduleGroup := group.Group("{{.RoutePath}}")
	moduleGroup.GET("/info", h.Info)
}

func (h *Handler) Info(c *gin.Context) {
	info, err := h.service.Info(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, info)
}
`

const handlerTestTemplate = `package {{.PackageName}}

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func TestModuleInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v2http.RegisterV2(engine, v2http.Config{}, NewModule(nil))

	req := httptest.NewRequest(http.MethodGet, "/v2{{.RoutePath}}/info", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body response.Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Code != response.Success || body.Message != "ok" {
		t.Fatalf("response = %+v, want success ok", body)
	}
	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("data type = %T, want map[string]any", body.Data)
	}
	if data["name"] != "{{.ModulePath}}" {
		t.Fatalf("name = %v, want {{.ModulePath}}", data["name"])
	}
}
`

const modelTemplate = `package {{.PackageName}}

type Info struct {
	Name string
	Message string
}
`

const moduleTemplate = `package {{.PackageName}}

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
)

type Module struct {
	handler *Handler
}

func NewModule(c *container.Container) *Module {
	_ = c
	repo := NewMemoryRepository()
	service := NewService(repo)
	return &Module{
		handler: NewHandler(service),
	}
}

func (m *Module) RegisterHTTP(routes v2http.Routes) {
	m.handler.Register(routes.Public)
}
`

const repositoryTemplate = `package {{.PackageName}}

import "context"

type Repository interface {
	Info(ctx context.Context) Info
}

type MemoryRepository struct{}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) Info(ctx context.Context) Info {
	return Info{
		Name: "{{.ModulePath}}",
		Message: "v2 module registered",
	}
}
`

const serviceTemplate = `package {{.PackageName}}

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Info(ctx context.Context) (InfoResponse, error) {
	info := s.repo.Info(ctx)
	return InfoResponse{
		Name: info.Name,
		Message: info.Message,
	}, nil
}
`
