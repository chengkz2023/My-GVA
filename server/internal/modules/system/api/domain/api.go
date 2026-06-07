package domain

type Api struct {
	ID          uint
	Path        string
	Description string
	ApiGroup    string
	Method      string
}

type SyncApi struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	ApiGroup    string `json:"apiGroup"`
	Description string `json:"description"`
}

type SyncResult struct {
	NewApis    []SyncApi `json:"newApis"`
	DeleteApis []SyncApi `json:"deleteApis"`
	IgnoreApis []SyncApi `json:"ignoreApis"`
}

type SyncRequest struct {
	NewApis    []SyncApi `json:"newApis"`
	DeleteApis []SyncApi `json:"deleteApis"`
	IgnoreApis []SyncApi `json:"ignoreApis"`
}

type IgnoreApiInput struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Flag   bool   `json:"flag"`
}

type RouteInfo struct {
	Path   string
	Method string
}
