package application

type FileResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	ClassID int    `json:"classId"`
	URL     string `json:"url"`
	Tag     string `json:"tag"`
	Key     string `json:"key"`
}

type ListResponse struct {
	List     []FileResponse `json:"list"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

type UploadResponse struct {
	File FileResponse `json:"file"`
}
