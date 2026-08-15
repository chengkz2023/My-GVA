package application

import (
	"context"
	"errors"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/api/domain"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/chengkz2023/My-GVA/server/internal/platform/authz"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
)

type Service struct {
	repo          domain.Repository
	policyProvider authz.PolicyProvider
	policySyncer  authz.PolicySyncer
	useStrictAuth bool
}

func NewService(repo domain.Repository, policyProvider authz.PolicyProvider, policySyncer authz.PolicySyncer, useStrictAuth bool) *Service {
	return &Service{repo: repo, policyProvider: policyProvider, policySyncer: policySyncer, useStrictAuth: useStrictAuth}
}

func (s *Service) List(ctx context.Context, query domain.ListQuery) (ListResponse, error) {
	if s.repo == nil {
		return ListResponse{List: []ApiResponse{}}, nil
	}
	result, err := s.repo.List(ctx, query)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return ListResponse{List: []ApiResponse{}}, nil
	}
	if err != nil {
		return ListResponse{}, apperrors.New(apperrors.Internal, 0, "list apis failed", err)
	}
	return ListResponse{
		List:     mapApis(result.List),
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}

func (s *Service) GetAll(ctx context.Context) (AllResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return AllResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return AllResponse{List: []ApiResponse{}}, nil
	}
	apis, err := s.repo.GetAll(ctx)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return AllResponse{List: []ApiResponse{}}, nil
	}
	if err != nil {
		return AllResponse{}, apperrors.New(apperrors.Internal, 0, "get all apis failed", err)
	}

	if s.useStrictAuth && s.policyProvider != nil {
		policies, _ := s.policyProvider.Policies(actor.AuthorityID)
		if len(policies) == 0 {
			// 严格模式下无策略的角色不应看到全部 API
			return AllResponse{List: []ApiResponse{}}, nil
		}
		policySet := make(map[string]map[string]bool, len(policies))
		for _, p := range policies {
			if policySet[p.Path] == nil {
				policySet[p.Path] = make(map[string]bool)
			}
			policySet[p.Path][p.Method] = true
		}
		filtered := make([]domain.Api, 0, len(apis))
		for _, api := range apis {
			if policySet[api.Path] != nil && policySet[api.Path][api.Method] {
				filtered = append(filtered, api)
			}
		}
		apis = filtered
	}
	return AllResponse{List: mapApis(apis)}, nil
}

func (s *Service) Groups(ctx context.Context) (GroupsResponse, error) {
	if s.repo == nil {
		return GroupsResponse{Groups: []string{}}, nil
	}
	groups, err := s.repo.Groups(ctx)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return GroupsResponse{Groups: []string{}}, nil
	}
	if err != nil {
		return GroupsResponse{}, apperrors.New(apperrors.Internal, 0, "get api groups failed", err)
	}
	if groups == nil {
		groups = []string{}
	}
	return GroupsResponse{Groups: groups}, nil
}

func (s *Service) Policies(ctx context.Context, authorityID uint) (PoliciesResponse, error) {
	if s.policyProvider == nil {
		return PoliciesResponse{Paths: []PolicyResponse{}}, nil
	}
	pols, err := s.policyProvider.Policies(authorityID)
	if err != nil {
		return PoliciesResponse{}, apperrors.New(apperrors.Internal, 0, "get policies failed", err)
	}
	paths := make([]PolicyResponse, 0, len(pols))
	for _, p := range pols {
		paths = append(paths, PolicyResponse{Path: p.Path, Method: p.Method})
	}
	return PoliciesResponse{Paths: paths}, nil
}

func (s *Service) UpdatePolicies(ctx context.Context, authorityID uint, policies []PolicyResponse) error {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.policySyncer == nil {
		return apperrors.WithMessage(apperrors.Internal, "policy syncer unavailable")
	}
	if s.policyProvider != nil && s.repo != nil && s.useStrictAuth {
		allApis, err := s.repo.GetAll(ctx)
		if err == nil {
			allowed := make(map[string]map[string]bool)
			for _, api := range allApis {
				if allowed[api.Path] == nil {
					allowed[api.Path] = make(map[string]bool)
				}
				allowed[api.Path][api.Method] = true
			}
			for _, p := range policies {
				if _, ok := allowed[p.Path]; !ok || !allowed[p.Path][p.Method] {
					return apperrors.WithMessage(apperrors.Forbidden, "api not in allowed scope")
				}
			}
		}
	}

	items := make([]authz.Policy, 0, len(policies))
	for _, p := range policies {
		items = append(items, authz.Policy{Path: p.Path, Method: p.Method})
	}
	if err := s.policySyncer.SyncPolicies(authorityID, items); err != nil {
		return apperrors.New(apperrors.Internal, 0, "sync policies failed", err)
	}
	return nil
}

func (s *Service) Create(ctx context.Context, input domain.SaveApiInput) (ApiResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return ApiResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return ApiResponse{}, apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	id, err := s.repo.Save(ctx, input)
	if errors.Is(err, domain.ErrApiDuplicate) {
		return ApiResponse{}, apperrors.WithMessage(apperrors.Conflict, "存在相同api")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return ApiResponse{}, apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	if err != nil {
		return ApiResponse{}, apperrors.New(apperrors.Internal, 0, "create api failed", err)
	}
	return ApiResponse{ID: id, Path: input.Path, Description: input.Description, ApiGroup: input.ApiGroup, Method: input.Method}, nil
}

func (s *Service) Update(ctx context.Context, input domain.SaveApiInput) (ApiResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return ApiResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return ApiResponse{}, apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	id, err := s.repo.Save(ctx, input)
	if errors.Is(err, domain.ErrApiDuplicate) {
		return ApiResponse{}, apperrors.WithMessage(apperrors.Conflict, "存在相同api路径")
	}
	if errors.Is(err, domain.ErrApiNotFound) {
		return ApiResponse{}, apperrors.WithMessage(apperrors.NotFound, "api not found")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return ApiResponse{}, apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	if err != nil {
		return ApiResponse{}, apperrors.New(apperrors.Internal, 0, "update api failed", err)
	}
	return ApiResponse{ID: id, Path: input.Path, Description: input.Description, ApiGroup: input.ApiGroup, Method: input.Method}, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	deleted, err := s.repo.Delete(ctx, id)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	if err != nil {
		return apperrors.New(apperrors.Internal, 0, "delete api failed", err)
	}
	if s.policySyncer != nil {
		_ = s.policySyncer.RemovePolicies(deleted.Path, deleted.Method)
	}
	return nil
}

func (s *Service) BatchDelete(ctx context.Context, ids []int) error {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	deleted, err := s.repo.DeleteByIds(ctx, ids)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	if err != nil {
		return apperrors.New(apperrors.Internal, 0, "batch delete apis failed", err)
	}
	if s.policySyncer != nil {
		for _, d := range deleted {
			_ = s.policySyncer.RemovePolicies(d.Path, d.Method)
		}
	}
	return nil
}

func (s *Service) FreshCasbin() error {
	if s.policySyncer == nil {
		return apperrors.WithMessage(apperrors.Internal, "policy syncer unavailable")
	}
	if err := s.policySyncer.RefreshPolicies(); err != nil {
		return apperrors.New(apperrors.Internal, 0, "refresh casbin failed", err)
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (ApiResponse, error) {
	if s.repo == nil {
		return ApiResponse{}, apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	api, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return ApiResponse{}, apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	if err != nil {
		return ApiResponse{}, apperrors.New(apperrors.Internal, 0, "get api failed", err)
	}
	return ApiResponse{ID: api.ID, Path: api.Path, Description: api.Description, ApiGroup: api.ApiGroup, Method: api.Method}, nil
}

func (s *Service) Sync(routes []domain.RouteInfo) (domain.SyncResult, error) {
	if s.repo == nil {
		return domain.SyncResult{}, nil
	}
	allApis, err := s.repo.GetAll(context.Background())
	if err != nil {
		return domain.SyncResult{}, err
	}
	ignored, err := s.repo.GetIgnored(context.Background())
	if err != nil {
		return domain.SyncResult{}, err
	}

	apiMap := make(map[string]map[string]bool)
	for _, a := range allApis {
		if apiMap[a.Path] == nil {
			apiMap[a.Path] = make(map[string]bool)
		}
		apiMap[a.Path][a.Method] = true
	}
	ignoredMap := make(map[string]map[string]bool)
	for _, ig := range ignored {
		if ignoredMap[ig.Path] == nil {
			ignoredMap[ig.Path] = make(map[string]bool)
		}
		ignoredMap[ig.Path][ig.Method] = true
	}

	newApis := make([]domain.SyncApi, 0)
	deleteApis := make([]domain.SyncApi, 0)
	ignoreApis := make([]domain.SyncApi, 0)

	routeMap := make(map[string]map[string]bool)
	for _, r := range routes {
		if routeMap[r.Path] == nil {
			routeMap[r.Path] = make(map[string]bool)
		}
		routeMap[r.Path][r.Method] = true
		if _, ok := apiMap[r.Path]; !ok || !apiMap[r.Path][r.Method] {
			if _, ok := ignoredMap[r.Path]; ok && ignoredMap[r.Path][r.Method] {
				ignoreApis = append(ignoreApis, domain.SyncApi{Path: r.Path, Method: r.Method})
			} else {
				newApis = append(newApis, domain.SyncApi{Path: r.Path, Method: r.Method})
			}
		}
	}
	for _, a := range allApis {
		if _, ok := routeMap[a.Path]; !ok || !routeMap[a.Path][a.Method] {
			deleteApis = append(deleteApis, domain.SyncApi{
				Path: a.Path, Method: a.Method, ApiGroup: a.ApiGroup, Description: a.Description,
			})
		}
	}
	return domain.SyncResult{
		NewApis:    newApis,
		DeleteApis: deleteApis,
		IgnoreApis: ignoreApis,
	}, nil
}

func (s *Service) IgnoreApi(ctx context.Context, input domain.IgnoreApiInput) error {
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	return s.repo.Ignore(ctx, input.Path, input.Method, input.Flag)
}

func (s *Service) BatchSync(ctx context.Context, req domain.SyncRequest) error {
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "api repository unavailable")
	}
	newInputs := make([]domain.SaveApiInput, 0, len(req.NewApis))
	for _, a := range req.NewApis {
		newInputs = append(newInputs, domain.SaveApiInput{
			Path: a.Path, Method: a.Method, ApiGroup: a.ApiGroup, Description: a.Description,
		})
	}
	deleteInputs := make([]domain.SaveApiInput, 0, len(req.DeleteApis))
	for _, a := range req.DeleteApis {
		deleteInputs = append(deleteInputs, domain.SaveApiInput{
			Path: a.Path, Method: a.Method,
		})
	}
	ignoreInputs := make([]domain.SaveApiInput, 0, len(req.IgnoreApis))
	for _, a := range req.IgnoreApis {
		ignoreInputs = append(ignoreInputs, domain.SaveApiInput{
			Path: a.Path, Method: a.Method,
		})
	}
	if err := s.repo.BatchCreate(ctx, newInputs); err != nil {
		return apperrors.New(apperrors.Internal, 0, "batch create apis failed", err)
	}
	if err := s.repo.BatchDeleteByPathMethod(ctx, deleteInputs); err != nil {
		return apperrors.New(apperrors.Internal, 0, "batch delete apis failed", err)
	}
	for _, ig := range ignoreInputs {
		_ = s.repo.Ignore(ctx, ig.Path, ig.Method, true)
	}
	return nil
}

func mapApis(apis []domain.Api) []ApiResponse {
	items := make([]ApiResponse, 0, len(apis))
	for _, api := range apis {
		items = append(items, ApiResponse{
			ID:          api.ID,
			Path:        api.Path,
			Description: api.Description,
			ApiGroup:    api.ApiGroup,
			Method:      api.Method,
		})
	}
	return items
}
