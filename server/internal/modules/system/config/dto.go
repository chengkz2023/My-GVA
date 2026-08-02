package config

import platformconfig "github.com/chengkz2023/My-GVA/server/internal/platform/config"

type InfoResponse struct {
	Config platformconfig.Snapshot `json:"config"`
}
