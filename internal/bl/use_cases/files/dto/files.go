package dto

import (
	"io"
	"time"
)

type File struct {
	Uuid      string     // Uuid файла
	UserUuid  string     // Пользователь, изменивший файл
	Name      string     // Название файла
	CheckSum  []byte     // Контрольная сумма файла
	StatusId  uint       // Статус файла
	TypeId    uint       // Тип файла
	Reader    io.Reader  // Данные файла при загрузке
	Data      []byte     // Данные файла при скачивании
	Size      int64      // Размер файла
	CreatedAt *time.Time // Дата создания файла
	UpdatedAt *time.Time // Дата обновления файла
}

type FileHFDTO struct {
	Uuid         string
	RedactorUuid string
	Type         string
	CheckSum     string
	Status       int
	History      []History
}

type History struct {
	RedactorUuid string
	Status       int
	UpdatedAt    string
}
