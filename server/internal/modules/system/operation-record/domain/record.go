package domain

type Record struct {
	ID           uint
	IP           string
	Method       string
	Path         string
	Status       int
	Latency      int64
	Agent        string
	ErrorMessage string
	Body         string
	Resp         string
	UserID       int
	Username     string
	NickName     string
}
