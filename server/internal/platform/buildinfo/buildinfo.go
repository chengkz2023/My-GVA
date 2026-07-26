package buildinfo

const (
	Version     = "v2.8.8"
	AppName     = "Gin-Vue-Admin"
	Description = "使用gin+vue进行极速开发的全栈开发基础平台"
)

type Info struct {
	AppName     string
	Version     string
	Description string
}

func Current() Info {
	return Info{
		AppName:     AppName,
		Version:     Version,
		Description: Description,
	}
}
