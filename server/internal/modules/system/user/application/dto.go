package application

import "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"

type CurrentUserResponse struct {
	User SourceUser `json:"user"`
}

type ListUsersQuery struct {
	Page     pagination.Page
	Username string
	NickName string
	Phone    string
	Email    string
}

type ListUsersResponse = pagination.Result[SourceUser]

type ChangePasswordCommand struct {
	OldPassword string
	NewPassword string
}

type ChangePasswordResponse struct {
	Changed bool `json:"changed"`
}

type UpdateProfileCommand struct {
	NickName  string
	HeaderImg string
	Phone     string
	Email     string
}

type UpdateProfileResponse struct {
	User SourceUser `json:"user"`
}

type CreateUserCommand struct {
	Username     string
	Password     string
	NickName     string
	HeaderImg    string
	AuthorityID  uint
	AuthorityIDs []uint
	Enable       int
	Phone        string
	Email        string
}

type CreateUserResponse struct {
	User SourceUser `json:"user"`
}

type DeleteUserResponse struct {
	DeletedID uint `json:"deletedId"`
}

type ResetPasswordCommand struct {
	UserID   uint
	Password string
}

type ResetPasswordResponse struct {
	Changed bool `json:"changed"`
}

type SetAuthoritiesCommand struct {
	UserID       uint
	AuthorityIDs []uint
}

type SetAuthoritiesResponse struct {
	User SourceUser `json:"user"`
}

type SourceUser struct {
	ID           uint   `json:"ID"`
	UUID         string `json:"uuid,omitempty"`
	Username     string `json:"username"`
	UserName     string `json:"userName"`
	NickName     string `json:"nickName"`
	HeaderImg    string `json:"headerImg,omitempty"`
	AuthorityID  uint   `json:"authorityId"`
	Authorities  []any  `json:"authorities"`
	AuthorityIds []uint `json:"authorityIds"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	Enable       int    `json:"enable,omitempty"`
	Source       string `json:"source"`
}
