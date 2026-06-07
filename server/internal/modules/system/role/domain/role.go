package domain

type Role struct {
	AuthorityID   uint
	AuthorityName string
	ParentID      *uint
	DefaultRouter string
	Children      []Role
}

type SaveRoleInput struct {
	AuthorityID   uint
	AuthorityName string
	ParentID      *uint
	DefaultRouter string
}

type DataAuthorityInput struct {
	AuthorityID     uint
	DataAuthorityIDs []uint
}
