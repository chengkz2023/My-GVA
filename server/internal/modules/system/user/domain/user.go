package domain

type User struct {
	ID          uint
	UUID        string
	Username    string
	NickName    string
	HeaderImg   string
	AuthorityID uint
	Phone       string
	Email       string
	Enable      int
}
