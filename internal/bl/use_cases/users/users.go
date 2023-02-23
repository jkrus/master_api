package users

import (
	"context"

	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/users/dto"
	"github.com/jkrus/master_api/pkg/errors"
)

type UsersI interface {
	Create(ctx context.Context, data *dto.User) (*dto.User, error)
	GetByUuid(ctx context.Context, uuid string) (*dto.User, error)
	Update(ctx context.Context, uuid string, data *dto.User) (*dto.User, error)
	Delete(ctx context.Context, uuid string) error
}

type users struct {
	di internal.IAppDeps
}

func NewUsersI(di internal.IAppDeps) UsersI {
	return &users{di: di}
}

func (u *users) Create(ctx context.Context, data *dto.User) (*dto.User, error) {

	created, err := u.di.DBRepo().UserRepository.Users.Create(ctx, data)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}
	return created, nil
}

func (u *users) GetByUuid(ctx context.Context, uuid string) (*dto.User, error) {
	user, err := u.di.DBRepo().UserRepository.Users.GetByUuid(ctx, uuid)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}

	return user, nil
}

func (u *users) Update(ctx context.Context, uuid string, data *dto.User) (*dto.User, error) {
	updated, err := u.di.DBRepo().UserRepository.Users.Update(ctx, uuid, data)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}

	return updated, nil
}

func (u *users) Delete(ctx context.Context, uuid string) error {
	return u.di.DBRepo().UserRepository.Users.Delete(ctx, uuid)
}
