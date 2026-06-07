package buildinfo

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Info struct {
	AppName     string
	Version     string
	Description string
}

func Current() Info {
	return Info{
		AppName:     global.AppName,
		Version:     global.Version,
		Description: global.Description,
	}
}
