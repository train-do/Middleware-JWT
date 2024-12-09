package service

import (
	"voucher_system/models"
	"voucher_system/repository"

	"go.uber.org/zap"
)

type UserService interface {
	Login(email string) (models.User, error)
	Register(user models.User) error
}

type userService struct {
	Repo repository.Repository
	log  *zap.Logger
}

func NewUserService(repo repository.Repository, log *zap.Logger) UserService {
	return &userService{Repo: repo, log: log}
}

func (s *userService) Login(email string) (models.User, error) {
	return s.Repo.User.Login(email)
}
func (s *userService) Register(user models.User) error {
	return s.Repo.User.Register(user)
}
