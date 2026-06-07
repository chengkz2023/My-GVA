package pagination

type Page struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

func Normalize(page Page) Page {
	if page.Page <= 0 {
		page.Page = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = 10
	}
	if page.PageSize > 100 {
		page.PageSize = 100
	}
	return page
}

func (p Page) Limit() int {
	return p.PageSize
}

func (p Page) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return p.PageSize * (p.Page - 1)
}

type Result[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}
