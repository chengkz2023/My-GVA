package application

type TreeResponse struct {
	List []RoleResponse `json:"list"`
}

type RoleResponse struct {
	AuthorityID   uint           `json:"authorityId"`
	AuthorityName string         `json:"authorityName"`
	ParentID      *uint          `json:"parentId"`
	DefaultRouter string         `json:"defaultRouter"`
	Children      []RoleResponse `json:"children"`
}

type RoleDetailResponse struct {
	AuthorityID      uint   `json:"authorityId"`
	AuthorityName    string `json:"authorityName"`
	ParentID         *uint  `json:"parentId"`
	DefaultRouter    string `json:"defaultRouter"`
	DataAuthorityIDs []uint `json:"dataAuthorityIds"`
}

type CopyRequest struct {
	OldAuthorityID uint                  `json:"oldAuthorityId"`
	Authority      CreateRoleRequest     `json:"authority"`
}

type CreateRoleRequest struct {
	AuthorityID   uint   `json:"authorityId"`
	AuthorityName string `json:"authorityName"`
	ParentID      *uint  `json:"parentId"`
	DefaultRouter string `json:"defaultRouter"`
}

type DataAuthorityRequest struct {
	AuthorityID      uint   `json:"authorityId"`
	DataAuthorityIDs []uint `json:"dataAuthorityIds"`
}
