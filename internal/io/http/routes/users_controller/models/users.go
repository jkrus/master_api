package models

import "github.com/jkrus/master_api/internal/bl/use_cases/users/dto"

// User Пользователь
type User struct {
	Uuid      string // Uuid модель
	PublicKey string // Публичный ключ
	Role      string // Роль
	Status    string // Статус
}

func (fs *User) ToDTO() *dto.User {
	if fs == nil {
		return nil
	}

	return &dto.User{
		Uuid:      fs.Uuid,
		PublicKey: fs.PublicKey,
		Role:      fs.Role,
		Status:    fs.Status,
	}
}

func (fs *User) FromDTO(model *dto.User) *User {
	if model == nil {
		return nil
	}

	return &User{
		Uuid:      model.Uuid,
		PublicKey: model.PublicKey,
		Role:      model.Role,
		Status:    model.Status,
	}
}
