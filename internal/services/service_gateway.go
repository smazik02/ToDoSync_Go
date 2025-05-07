package services

import (
	"todosync_go/internal/shared"
	"todosync_go/utils"
)

type ServiceGateway struct {
	userService *UserService
}

func NewServiceGateway(userService *UserService) *ServiceGateway {
	return &ServiceGateway{
		userService: userService,
	}
}

func (g ServiceGateway) Direct(resourceMethod utils.ResourceMethod, payload []byte, client *shared.Client) (*ServiceResponse, error) {
	_, isLoggedIn := g.userService.loggedInUsers[client.UserId]
	if (!isLoggedIn || !client.IsLoggedIn()) && resourceMethod != utils.AUTH_LOGIN {
		return nil, ServiceError{"Not logged in", AUTH}
	}

	switch resourceMethod {
	case utils.AUTH_LOGIN:
		return g.userService.LoginUser(payload, client)
	default:
		return nil, ServiceError{"Unknown operation", AUTH}
	}
}
