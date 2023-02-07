package dto

import (
	"io"
)

type FileIN struct {
	UUID   string // UUID файла
	Name   string // Имя файла
	Size   int64  // Размер файла
	Reader io.Reader
}

type FileOUT struct {
	UUID  string // UUID файла
	Name  string // Имя файла
	Size  int64  // Размер файла
	Bytes []byte // Данные файла
}
