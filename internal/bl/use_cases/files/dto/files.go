package dto

import (
	"io"
)

type File struct {
	Uuid     string    // Uuid файла
	UserUuid string    // Пользователь, изменивший файл
	Name     string    // Название файла
	CheckSum []byte    // Контрольная сумма файла
	StatusId uint      // Статус файла
	TypeId   uint      // Тип файла
	Reader   io.Reader // Данные файла при загрузке
	Data     []byte    // Данные файла при скачивании
	Size     int64     // Размер файла
}
