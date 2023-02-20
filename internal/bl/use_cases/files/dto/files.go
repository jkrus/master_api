package dto

import (
	"io"
)

type FileIN struct {
	Uuid   string // Uuid файла
	Name   string // Имя файла
	Size   uint   // Размер файла
	Reader io.Reader
}

type FileOUT struct {
	Uuid  string // Uuid файла
	Name  string // Имя файла
	Size  uint   // Размер файла
	Bytes []byte // Данные файла
}

type FileINHF struct {
	Uuid         string
	RedactorUuid string
	Type         string
	CheckSum     string
	Status       int
}
