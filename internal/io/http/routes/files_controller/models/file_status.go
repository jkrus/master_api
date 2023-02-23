package models

import "github.com/jkrus/master_api/internal/bl/use_cases/files/dto"

type FileStatus struct {
	Id          uint   // Id статуса файла
	Status      string // Название статуса файла
	Description string // Описание статуса файла
}

func (fs *FileStatus) ToDTO() *dto.FileStatus {
	if fs == nil {
		return nil
	}

	return &dto.FileStatus{
		Id:          fs.Id,
		Status:      fs.Status,
		Description: fs.Description,
	}
}

func (fs *FileStatus) FromDTO(model *dto.FileStatus) *FileStatus {
	if model == nil {
		return nil
	}

	return &FileStatus{
		Id:          model.Id,
		Status:      model.Status,
		Description: model.Description,
	}
}
