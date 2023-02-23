package files

import (
	"context"

	"gorm.io/gorm"

	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/internal/common/err_const"
	"github.com/jkrus/master_api/internal/stores/db/repo/base"
)

// FileStatus Статус файла
type FileStatus struct {
	base.DictionaryIDModel        // Uuid модель
	Status                 string `gorm:"index"` // Название статуса файла
	Description            string // Описание статуса файла
}

func (fs *FileStatus) toDTO() *dto.FileStatus {
	if fs == nil {
		return nil
	}

	return &dto.FileStatus{
		Id:          fs.ID,
		Status:      fs.Status,
		Description: fs.Description,
	}
}

func (fs *FileStatus) fromDTO(v *dto.FileStatus) {
	if v == nil {
		return
	}

	fs.ID = v.Id
	fs.Status = v.Status
	fs.Description = v.Description

	return
}

type FileStatusRepositoryI interface {
	Create(ctx context.Context, data *dto.FileStatus) (*dto.FileStatus, error)
	GetById(ctx context.Context, fileStatusId uint) (*dto.FileStatus, error)
	Update(ctx context.Context, fileStatusId uint, data *dto.FileStatus) (*dto.FileStatus, error)
	Delete(ctx context.Context, fileStatusId uint) error
}

type fileStatusRepository struct {
	db *gorm.DB
}

func NewFileStatusRepository(db *gorm.DB) FileStatusRepositoryI {
	return &fileStatusRepository{
		db: db,
	}
}

func (f *fileStatusRepository) Create(ctx context.Context, data *dto.FileStatus) (*dto.FileStatus, error) {
	result := &FileStatus{}
	result.fromDTO(data)

	err := f.db.WithContext(ctx).Create(result).Error
	if err != nil {
		return nil, err
	}

	return result.toDTO(), nil
}

func (f *fileStatusRepository) GetById(ctx context.Context, fileStatusId uint) (*dto.FileStatus, error) {
	result := &FileStatus{}

	err := f.db.WithContext(ctx).Find(&result, fileStatusId).Error
	if err != nil {
		return nil, err
	}

	return result.toDTO(), nil
}

func (f *fileStatusRepository) Update(ctx context.Context, fileStatusId uint, data *dto.FileStatus) (*dto.FileStatus, error) {
	update := &FileStatus{}
	update.fromDTO(data)
	update.ID = fileStatusId

	tx := f.db.WithContext(ctx).Save(update)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, err_const.ErrDatabaseRecordNotFound
	}

	return update.toDTO(), nil
}

func (f *fileStatusRepository) Delete(ctx context.Context, fileStatusId uint) error {
	return f.db.WithContext(ctx).Delete(&FileStatus{}, fileStatusId).Error
}
