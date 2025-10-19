package response

import "ApkAdmin/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
