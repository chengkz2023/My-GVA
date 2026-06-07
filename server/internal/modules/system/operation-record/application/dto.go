package application

type RecordResponse struct {
	ID           uint   `json:"id"`
	IP           string `json:"ip"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	Status       int    `json:"status"`
	Latency      int64  `json:"latency"`
	Agent        string `json:"agent"`
	ErrorMessage string `json:"errorMessage"`
	Body         string `json:"body"`
	Resp         string `json:"resp"`
	UserID       int    `json:"userId"`
	Username     string `json:"username"`
	NickName     string `json:"nickName"`
}

type ListResponse struct {
	List     []RecordResponse `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}
