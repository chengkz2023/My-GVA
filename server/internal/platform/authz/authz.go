package authz

import (
	"context"

	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
)

type Authorizer interface {
	Can(ctx context.Context, actor platformauth.Actor, action string, resource any) error
}

type Policy struct {
	Path   string
	Method string
}

type PolicyProvider interface {
	Policies(authorityID uint) ([]Policy, error)
}

type PolicySyncer interface {
	SyncPolicies(authorityID uint, policies []Policy) error
	RemovePolicies(path, method string) error
	RefreshPolicies() error
}

type AuthorityChecker interface {
	CheckAuthorityAuth(ctx context.Context, adminID, targetID uint) error
	GetDescendantIDs(ctx context.Context, authorityID uint) ([]uint, error)
}
