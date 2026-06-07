package version

type InfoResponse struct {
	AppName     string `json:"appName"`
	Version     string `json:"version"`
	Description string `json:"description"`
}
