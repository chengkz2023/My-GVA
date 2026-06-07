package auth

import "context"

type Actor struct {
	UserID      uint   `json:"userId"`
	AuthorityID uint   `json:"authorityId"`
	Username    string `json:"username"`
	NickName    string `json:"nickName"`
}

type Authorizer interface {
	Can(ctx context.Context, actor Actor, action string, resource any) error
}

type actorContextKey struct{}

func ContextWithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, actorContextKey{}, actor)
}

func ActorFromContext(ctx context.Context) (Actor, bool) {
	actor, ok := ctx.Value(actorContextKey{}).(Actor)
	return actor, ok
}
