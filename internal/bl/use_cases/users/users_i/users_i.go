package users_i

import (
	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/users"
)

type UserLogic struct {
	Users users.UsersI
}

func NewUserLogic(di internal.IAppDeps) *UserLogic {
	return &UserLogic{
		Users: users.NewUsersI(di),
	}
}
