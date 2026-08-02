package buildinfo

const (
	Version     = "v1.0.0"
	AppName     = "BoyKing Admin"
	Description = "模块化后台管理脚手架"
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
