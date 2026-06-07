package config

import platformconfig "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/config"

type InfoResponse struct {
	Config platformconfig.Snapshot `json:"config"`
}
