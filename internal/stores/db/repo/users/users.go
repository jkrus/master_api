package users

import (
	"context"

	"github.com/jkrus/master_api/internal/bl/use_cases/users/dto"
	"github.com/jkrus/master_api/internal/common/err_const"

	"github.com/jkrus/master_api/internal/stores/db/repo/base"

	"gorm.io/gorm"
)

// User Пользователь
type User struct {
	base.UuidModel // Uuid модель

	PublicKey string // Публичный ключ
	Role      string // Роль
	Status    string // Статус
}

func (u *User) toDTO() *dto.User {
	if u == nil {
		return nil
	}

	return &dto.User{
		Uuid:      u.Uuid,
		PublicKey: u.PublicKey,
		Role:      u.Role,
		Status:    u.Status,
	}
}

func (u *User) fromDTO(v *dto.User) {
	if v == nil {
		return
	}

	u.UuidModel = base.UuidModel{Uuid: v.Uuid}
	u.PublicKey = v.PublicKey
	u.Role = v.Role
	u.Status = v.Status
}

type IUserRepository interface {
	Create(ctx context.Context, dtm *dto.User) (*dto.User, error)
	GetByUuid(ctx context.Context, uuid string) (*dto.User, error)
	Update(ctx context.Context, uuid string, data *dto.User) (*dto.User, error)
	Delete(ctx context.Context, uuid string) error

	WithTx(tx *gorm.DB) IUserRepository
}

type userRepository struct {
	db *gorm.DB
}

func NewUsersRepository(dbHandler *gorm.DB) IUserRepository {
	return &userRepository{db: dbHandler}
}

func (u *userRepository) Create(ctx context.Context, data *dto.User) (*dto.User, error) {
	result := &User{}
	result.fromDTO(data)

	err := u.db.WithContext(ctx).Create(result).Error
	if err != nil {
		return nil, err
	}

	return result.toDTO(), nil
}

func (u *userRepository) GetByUuid(ctx context.Context, uuid string) (*dto.User, error) {
	result := &User{}

	err := u.db.WithContext(ctx).Find(&result, "uuid = ?", uuid).Error
	if err != nil {
		return nil, err
	}

	return result.toDTO(), nil
}

func (u *userRepository) Update(ctx context.Context, uuid string, data *dto.User) (*dto.User, error) {
	update := &User{}
	update.fromDTO(data)
	update.Uuid = uuid

	tx := u.db.WithContext(ctx).Save(update)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, err_const.ErrDatabaseRecordNotFound
	}

	return update.toDTO(), nil
}

func (u *userRepository) Delete(ctx context.Context, uuid string) error {
	return u.db.WithContext(ctx).Delete(&User{}, "uuid = ?", uuid).Error
}

func (u *userRepository) WithTx(tx *gorm.DB) IUserRepository {
	return &userRepository{db: tx}
}
