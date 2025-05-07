package services

import (
	"database/sql"
	"encoding/json"
	"todosync_go/internal/repositories"
	"todosync_go/internal/shared"
	"todosync_go/utils"

	"github.com/go-playground/validator/v10"
)

type RegisterUser struct {
	Username string `validate:"required" json:"username"`
}

type UserService struct {
	repository    repositories.UserRepository
	validator     *validator.Validate
	loggedInUsers map[int]any
}

func NewUserService(db *sql.DB) UserService {
	return UserService{
		repository:    repositories.NewUserRepository(db),
		validator:     validator.New(validator.WithRequiredStructEnabled()),
		loggedInUsers: make(map[int]any),
	}
}

func (s *UserService) LoginUser(payload []byte, client *shared.Client) (*ServiceResponse, error) {
	registerUser := RegisterUser{}
	if err := json.Unmarshal(payload, &registerUser); err != nil {
		return nil, ServiceError{err.Error(), AUTH}
	}

	if err := s.validator.Struct(registerUser); err != nil {
		return nil, ServiceError{
			message: utils.AnalyzeStructError(err, AUTH),
			source:  AUTH,
		}
	}

	isTaken, err := s.repository.IsUsernameTaken(registerUser.Username)
	if err != nil {
		return nil, ServiceError{err.Error(), AUTH}
	}
	if isTaken {
		return nil, ServiceError{"User with that name already exists", AUTH}
	}

	userId, err := s.repository.AddUser(registerUser.Username)
	if err != nil {
		return nil, ServiceError{err.Error(), AUTH}
	}
	s.loggedInUsers[userId] = struct{}{}
	client.UserId = userId

	return &ServiceResponse{
		Message:      []byte("OK\n{}\n\n"),
		Notification: nil,
	}, nil
}
