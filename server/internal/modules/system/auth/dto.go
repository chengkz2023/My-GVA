package auth

type MeResponse struct {
	Actor ActorResponse `json:"actor"`
}

type ActorResponse struct {
	UserID      uint   `json:"userId"`
	AuthorityID uint   `json:"authorityId"`
	Username    string `json:"username"`
	NickName    string `json:"nickName"`
}

type LoginResponse struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}

type UserDTO struct {
	ID          uint   `json:"id"`
	UUID        string `json:"uuid"`
	Username    string `json:"username"`
	NickName    string `json:"nickName"`
	AuthorityID uint   `json:"authorityId"`
	HeaderImg   string `json:"headerImg,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
