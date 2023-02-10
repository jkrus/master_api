package models

import (
	"io"

	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
)

type File struct {
	UUID   string
	Name   string // Имя файла
	Size   uint   // Размер файла
	Reader io.Reader
}

func (f *File) FromDTO(model *dto.FileIN) *File {
	if model == nil {
		return nil
	}

	return &File{
		UUID: f.UUID,
		Name: model.Name,
		Size: model.Size,
	}
}

func (f *File) ToDTO() *dto.FileIN {
	if f == nil {
		return nil
	}

	return &dto.FileIN{
		Uuid:   f.UUID,
		Name:   f.Name,
		Size:   f.Size,
		Reader: f.Reader,
	}
}

type Reference struct {
	Value string
}
