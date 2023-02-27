package files

import (
	"context"

	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/internal/common/err_const"
	"github.com/jkrus/master_api/internal/stores/db/repo/base"

	"gorm.io/gorm"
)

// File Файл
type File struct {
	base.UuidModel        // Uuid модель
	UserUuid       string `gorm:"foreignKey:users TYPE:uuid"` // Пользователь, изменивший файл
	Name           string // Имя файла
	CheckSum       []byte // Контрольная сумма файла
	FileStatusID   uint   `gorm:"TYPE:integer references file_statuses"` // Статус файла
	FileTypeId     uint   `gorm:"TYPE:integer REFERENCES file_types"`    // Тип файла
}

func (f *File) toDTO() *dto.File {
	if f == nil {
		return nil
	}

	return &dto.File{
		Uuid:      f.Uuid,
		UserUuid:  f.UserUuid,
		CheckSum:  f.CheckSum,
		StatusId:  f.FileStatusID,
		TypeId:    f.FileTypeId,
		CreatedAt: &f.CreatedAt,
		UpdatedAt: &f.UpdatedAt,
	}
}

func (f *File) fromDTO(v *dto.File) {
	if v == nil {
		return
	}

	f.UuidModel.Uuid = v.Uuid
	f.UserUuid = v.UserUuid
	f.Name = v.Name
	f.CheckSum = v.CheckSum
	f.FileStatusID = v.StatusId
	f.FileTypeId = v.TypeId

}

type IFileRepository interface {
	Create(ctx context.Context, data *dto.File) (*dto.File, error)
	GetById(ctx context.Context, fileUuid string) (*dto.File, error)
	Delete(ctx context.Context, uuid string) error
	Update(ctx context.Context, uuid string, data *dto.File) (*dto.File, error)

	WithTx(tx *gorm.DB) IFileRepository
}

type fileRepository struct {
	db *gorm.DB
}

func (f *fileRepository) Create(ctx context.Context, data *dto.File) (*dto.File, error) {
	result := &File{}
	result.fromDTO(data)

	err := f.db.WithContext(ctx).Save(result).Error
	if err != nil {
		return nil, err
	}

	return result.toDTO(), nil
}

func (f *fileRepository) GetById(ctx context.Context, fileUuid string) (*dto.File, error) {
	result := &File{}

	err := f.db.WithContext(ctx).Find(&result, "uuid = ?", fileUuid).Error
	if err != nil {
		return nil, err
	}

	return result.toDTO(), nil
}

func (f *fileRepository) Update(ctx context.Context, uuid string, data *dto.File) (*dto.File, error) {
	update := &File{}
	update.fromDTO(data)
	update.Uuid = uuid

	tx := f.db.WithContext(ctx).Save(update)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, err_const.ErrDatabaseRecordNotFound
	}

	return update.toDTO(), nil
}

func (f *fileRepository) Delete(ctx context.Context, uuid string) error {
	return f.db.WithContext(ctx).Delete(&FileType{}, uuid).Error
}

func NewFileRepository(dbHandler *gorm.DB) IFileRepository {
	return &fileRepository{db: dbHandler}
}

func (f *fileRepository) WithTx(tx *gorm.DB) IFileRepository {
	return &fileRepository{db: tx}
}
