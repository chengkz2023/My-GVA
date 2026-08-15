package domain

type User struct {
	ID           uint
	UUID         string
	Username     string
	NickName     string
	HeaderImg    string
	AuthorityID  uint
	AuthorityIDs []uint
	Phone        string
	Email        string
	Enable       int
}

// AdminUpdateInput 管理员视角的用户资料/状态更新（区别于用户自改 ProfilePatch）。
type AdminUpdateInput struct {
	NickName  string
	HeaderImg string
	Phone     string
	Email     string
	Enable    int
}
