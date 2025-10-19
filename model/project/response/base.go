package response

import (
	"ApkAdmin/model/project"
)

type LoginResponse struct {
	User      project.User `json:"user"`
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expiresAt"`
}
