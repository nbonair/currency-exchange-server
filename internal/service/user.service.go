package service

import (
	"context"

	"github.com/nbonair/currency-exchange-server/internal/model"
	"github.com/nbonair/currency-exchange-server/internal/repo"
)

type UserService interface {
	Register(ctx context.Context, in *model.RegisterInput) (result int, err error)
	Login(ctx context.Context) error
	VerifyOTP(ctx context.Context) error
	UpdatePasswordRegister(ctx context.Context) error
}

type userService struct {
	repoUser repo.UserRepository
}

func NewUserService(repoUser repo.UserRepository) UserService {
	return &userService{
		repoUser: repoUser,
	}
}

// Login implements UserService.
func (u *userService) Login(ctx context.Context) error {
	panic("unimplemented")
}

// Register implements UserService.
func (u *userService) Register(ctx context.Context, in *model.RegisterInput) (result int, err error) {
	panic("unimplemented")
}

// UpdatePasswordRegister implements UserService.
func (u *userService) UpdatePasswordRegister(ctx context.Context) error {
	panic("unimplemented")
}

// VerifyOTP implements UserService.
func (u *userService) VerifyOTP(ctx context.Context) error {
	panic("unimplemented")
}
