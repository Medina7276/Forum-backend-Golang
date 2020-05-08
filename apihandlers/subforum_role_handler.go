package apihandlers

import "git.01.alem.school/qjawko/forum/service"

type SubforumRoleHandler struct {
	SubforumRoleService *service.SubforumRoleService
	Endpoint            string
}
