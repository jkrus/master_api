package files

import (
	"context"

	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/internal/stores/db/repo/base"

	"gorm.io/gorm"
)

// File Уведомление
type File struct {
	base.UuidModel // Uuid модель

	Name string
	Size uint
}

func (f *File) toDTO() *dto.FileOUT {
	if f == nil {
		return nil
	}

	return &dto.FileOUT{
		Uuid: f.Uuid,
		Name: f.Name,
		Size: f.Size,
	}
}

func (f *File) fromDTO(v *dto.FileIN) *File {
	if v == nil {
		return nil
	}

	return &File{
		UuidModel: base.UuidModel{Uuid: v.Uuid},
		Name:      v.Name,
		Size:      v.Size,
	}
}

type IFileRepository interface {
	GetByUuid(ctx context.Context, uuid string) (*dto.FileOUT, error)
	Create(ctx context.Context, dtm *dto.FileIN) (*dto.FileOUT, error)
	Delete(ctx context.Context, uuid string) error

	WithTx(tx *gorm.DB) IFileRepository
}

type fileRepository struct {
	db *gorm.DB
}

func (f *fileRepository) GetByUuid(ctx context.Context, uuid string) (*dto.FileOUT, error) {
	// TODO implement me
	panic("implement me")
}

func (f *fileRepository) Create(ctx context.Context, dtm *dto.FileIN) (*dto.FileOUT, error) {
	// TODO implement me
	panic("implement me")
}

func (f *fileRepository) Delete(ctx context.Context, uuid string) error {
	// TODO implement me
	panic("implement me")
}

func NewFileRepository(dbHandler *gorm.DB) IFileRepository {
	return &fileRepository{db: dbHandler}
}

func (n *fileRepository) WithTx(tx *gorm.DB) IFileRepository {
	return &fileRepository{db: tx}
}
