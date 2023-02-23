package models

import "github.com/jkrus/master_api/internal/bl/use_cases/files/dto"

type FileType struct {
	Id          uint   // Id статуса файла
	Type        string // Название типа файла
	Description string // Описание типа файла
}

func (fs *FileType) ToDTO() *dto.FileType {
	if fs == nil {
		return nil
	}

	return &dto.FileType{
		Id:          fs.Id,
		Type:        fs.Type,
		Description: fs.Description,
	}
}

func (fs *FileType) FromDTO(model *dto.FileType) *FileType {
	if model == nil {
		return nil
	}

	return &FileType{
		Id:          model.Id,
		Type:        model.Type,
		Description: model.Description,
	}
}
