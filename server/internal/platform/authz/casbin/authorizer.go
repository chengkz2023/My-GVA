package casbin

import (
	"context"
	"errors"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/chengkz2023/My-GVA/server/internal/platform/authz"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
)

var (
	_ authz.Authorizer     = (*Authorizer)(nil)
	_ authz.PolicyProvider = (*Authorizer)(nil)
	_ authz.PolicySyncer   = (*Authorizer)(nil)
)

type Authorizer struct {
	enforcer *casbin.SyncedCachedEnforcer
}

func NewAuthorizer(enforcer *casbin.SyncedCachedEnforcer) *Authorizer {
	return &Authorizer{enforcer: enforcer}
}

func (a *Authorizer) Enforce(sub, obj, act string) (bool, error) {
	if a.enforcer == nil {
		return false, errors.New("casbin unavailable")
	}
	return a.enforcer.Enforce(sub, obj, act)
}

func (a *Authorizer) Can(ctx context.Context, actor auth.Actor, action string, resource any) error {
	path, ok := resource.(string)
	if !ok {
		return apperrors.WithMessage(apperrors.Forbidden, "invalid resource")
	}
	if a.enforcer == nil {
		return apperrors.WithMessage(apperrors.Forbidden, "casbin unavailable")
	}
	sub := strconv.FormatUint(uint64(actor.AuthorityID), 10)
	allowed, err := a.enforcer.Enforce(sub, path, action)
	if err != nil {
		return apperrors.New(apperrors.Internal, 0, "casbin enforce failed", err)
	}
	if !allowed {
		return apperrors.WithMessage(apperrors.Forbidden, "permission denied")
	}
	return nil
}

func (a *Authorizer) Policies(authorityID uint) ([]authz.Policy, error) {
	if a.enforcer == nil {
		return nil, nil
	}
	rules, _ := a.enforcer.GetFilteredPolicy(0, strconv.FormatUint(uint64(authorityID), 10))
	policies := make([]authz.Policy, 0, len(rules))
	for _, rule := range rules {
		if len(rule) >= 3 {
			policies = append(policies, authz.Policy{Path: rule[1], Method: rule[2]})
		}
	}
	return policies, nil
}

func (a *Authorizer) SyncPolicies(authorityID uint, policies []authz.Policy) error {
	if a.enforcer == nil {
		return nil
	}
	authStr := strconv.FormatUint(uint64(authorityID), 10)
	_, _ = a.enforcer.RemoveFilteredPolicy(0, authStr)
	if len(policies) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(policies))
	rules := make([][]string, 0, len(policies))
	for _, p := range policies {
		key := authStr + "\x00" + p.Path + "\x00" + p.Method
		if seen[key] {
			continue
		}
		seen[key] = true
		rules = append(rules, []string{authStr, p.Path, p.Method})
	}
	_, _ = a.enforcer.AddPolicies(rules)
	return nil
}

func (a *Authorizer) RemovePolicies(path, method string) error {
	if a.enforcer == nil {
		return nil
	}
	_, _ = a.enforcer.RemoveFilteredPolicy(1, path, method)
	return nil
}

func (a *Authorizer) RefreshPolicies() error {
	if a.enforcer == nil {
		return nil
	}
	return a.enforcer.LoadPolicy()
}
