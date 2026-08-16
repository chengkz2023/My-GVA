package application

type GreetingResponse struct {
	ID        uint   `json:"ID"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	DeptID    uint   `json:"deptId"`
	CreatedAt string `json:"createdAt"`
}

type ListGreetingsResponse struct {
	List []GreetingResponse `json:"list"`
}

type CreateGreetingCommand struct {
	Message string
	Author  string
}

type CreateGreetingResponse struct {
	Greeting GreetingResponse `json:"greeting"`
}
