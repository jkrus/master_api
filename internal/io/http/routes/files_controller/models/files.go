package models

import (
	"io"

	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
)

type File struct {
	Uuid     string    // Uuid файла
	UserUuid string    // Пользователь, изменивший файл
	Name     string    // Название файла
	CheckSum []byte    `json:"-"` // Контрольная сумма файла
	StatusId uint      // Статус файла
	TypeId   uint      // Тип файла
	Reader   io.Reader `json:"-"` // Данные файла при загрузке
	Data     []byte    `json:"-"` // Данные файла при скачивании
	Size     int64     `json:"-"` // Размер файла
}

func (f *File) FromDTO(model *dto.File) *File {
	if model == nil {
		return nil
	}

	return &File{
		Uuid:     model.Uuid,
		UserUuid: model.UserUuid,
		Name:     model.Name,
		CheckSum: model.CheckSum,
		StatusId: model.StatusId,
		TypeId:   model.TypeId,
		Reader:   model.Reader,
		Data:     model.Data,
		Size:     f.Size,
	}
}

func (f *File) ToDTO() *dto.File {
	if f == nil {
		return nil
	}

	return &dto.File{
		Uuid:     f.Uuid,
		UserUuid: f.UserUuid,
		Name:     f.Name,
		CheckSum: f.CheckSum,
		StatusId: f.StatusId,
		TypeId:   f.TypeId,
		Reader:   f.Reader,
		Data:     f.Data,
		Size:     f.Size,
	}
}
