package auth

type MeResponse struct {
	UserInfo UserInfoResponse `json:"userInfo"`
}

type UserInfoResponse struct {
	ID            uint           `json:"ID"`
	UUID          string         `json:"uuid"`
	NickName      string         `json:"nickName"`
	HeaderImg     string         `json:"headerImg"`
	Authority     AuthorityInfo  `json:"authority"`
	Authorities   []any          `json:"authorities"`
	OriginSetting map[string]any `json:"originSetting,omitempty"`
}

type AuthorityInfo struct {
	DefaultRouter string `json:"defaultRouter"`
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
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

type CaptchaResponse struct {
	CaptchaId     string `json:"captchaId"`
	PicPath       string `json:"picPath"`
	CaptchaLength int    `json:"captchaLength"`
	OpenCaptcha   bool   `json:"openCaptcha"`
}
