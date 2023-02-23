package users_i

import (
	"gorm.io/gorm"

	"github.com/jkrus/master_api/internal/stores/db/repo/users"
)

type UsersDBStore struct {
	Users users.IUserRepository
}

func NewUsersDBStore(dbHandler *gorm.DB) *UsersDBStore {
	return &UsersDBStore{
		Users: users.NewUsersRepository(dbHandler),
	}
}
