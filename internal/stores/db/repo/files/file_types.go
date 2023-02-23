package files

import (
	"context"

	"gorm.io/gorm"

	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/internal/common/err_const"
	"github.com/jkrus/master_api/internal/stores/db/repo/base"
)

// FileType Статус файла
type FileType struct {
	base.DictionaryIDModel        // Uuid модель
	Type                   string `gorm:"index"` // Название типа файла
	Description            string // Описание типа файла
}

func (fs *FileType) toDTO() *dto.FileType {
	if fs == nil {
		return nil
	}

	return &dto.FileType{
		Id:          fs.ID,
		Type:        fs.Type,
		Description: fs.Description,
	}
}

func (fs *FileType) fromDTO(v *dto.FileType) {
	if v == nil {
		return
	}

	fs.ID = v.Id
	fs.Type = v.Type
	fs.Description = v.Description

	return
}

type FileTypeRepositoryI interface {
	Create(ctx context.Context, data *dto.FileType) (*dto.FileType, error)
	GetById(ctx context.Context, FileTypeId uint) (*dto.FileType, error)
	Update(ctx context.Context, FileTypeId uint, data *dto.FileType) (*dto.FileType, error)
	Delete(ctx context.Context, FileTypeId uint) error
}

type FileTypeRepository struct {
	db *gorm.DB
}

func NewFileTypeRepository(db *gorm.DB) FileTypeRepositoryI {
	return &FileTypeRepository{
		db: db,
	}
}

func (f *FileTypeRepository) Create(ctx context.Context, data *dto.FileType) (*dto.FileType, error) {
	result := &FileType{}
	result.fromDTO(data)

	err := f.db.WithContext(ctx).Create(result).Error
	if err != nil {
		return nil, err
	}

	return result.toDTO(), nil
}

func (f *FileTypeRepository) GetById(ctx context.Context, FileTypeId uint) (*dto.FileType, error) {
	result := &FileType{}

	err := f.db.WithContext(ctx).Find(&result, FileTypeId).Error
	if err != nil {
		return nil, err
	}

	return result.toDTO(), nil
}

func (f *FileTypeRepository) Update(ctx context.Context, FileTypeId uint, data *dto.FileType) (*dto.FileType, error) {
	update := &FileType{}
	update.fromDTO(data)
	update.ID = FileTypeId

	tx := f.db.WithContext(ctx).Save(update)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, err_const.ErrDatabaseRecordNotFound
	}

	return update.toDTO(), nil
}

func (f *FileTypeRepository) Delete(ctx context.Context, FileTypeId uint) error {
	return f.db.WithContext(ctx).Delete(&FileType{}, FileTypeId).Error
}
