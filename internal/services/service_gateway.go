package services

import (
	"todosync_go/internal/parser"
	"todosync_go/internal/shared"
)

type ServiceGateway struct {
	userService *UserService
}

func NewServiceGateway(userService *UserService) *ServiceGateway {
	return &ServiceGateway{
		userService: userService,
	}
}

func (g ServiceGateway) Direct(parsedMessage *parser.ParserOutput, client *shared.Client) (*ServiceResponse, error) {
	_, isLoggedIn := g.userService.loggedInUsers[client.UserId]
	if (!isLoggedIn || !client.IsLoggedIn()) && parsedMessage.ResourceMethod != parser.AuthLogin {
		return nil, ServiceError{"Not logged in", AUTH}
	}

	switch parsedMessage.ResourceMethod {
	case parser.AuthLogin:
		return g.userService.LoginUser(parsedMessage.Payload, client)
	default:
		return nil, ServiceError{"Unknown operation", AUTH}
	}
}
