package application

type ApiResponse struct {
	ID          uint   `json:"id"`
	Path        string `json:"path"`
	Description string `json:"description"`
	ApiGroup    string `json:"apiGroup"`
	Method      string `json:"method"`
}

type ListResponse struct {
	List     []ApiResponse `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

type AllResponse struct {
	List []ApiResponse `json:"list"`
}

type GroupsResponse struct {
	Groups []string `json:"groups"`
}

type PolicyResponse struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type PoliciesResponse struct {
	Paths []PolicyResponse `json:"paths"`
}
