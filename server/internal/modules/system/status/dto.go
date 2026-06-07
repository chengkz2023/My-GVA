package status

type InfoResponse struct {
	Status   string         `json:"status"`
	Checks   ChecksResponse `json:"checks"`
	Warnings []string       `json:"warnings,omitempty"`
}

type ChecksResponse struct {
	Database DependencyStatus `json:"database"`
}

type DependencyStatus struct {
	Configured bool   `json:"configured"`
	OK         bool   `json:"ok"`
	Message    string `json:"message,omitempty"`
}
