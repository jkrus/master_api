package files

import (
	"context"

	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/pkg/errors"
)

type FileTypesI interface {
	Create(ctx context.Context, data *dto.FileType) (*dto.FileType, error)
	GetById(ctx context.Context, fileTypeId uint) (*dto.FileType, error)
	Update(ctx context.Context, fileTypeId uint, data *dto.FileType) (*dto.FileType, error)
	Delete(ctx context.Context, fileTypeId uint) error
}

type fileType struct {
	di internal.IAppDeps
}

func NewFileTypeI(di internal.IAppDeps) FileTypesI {
	return &fileType{
		di: di,
	}
}

func (fs *fileType) Create(ctx context.Context, data *dto.FileType) (*dto.FileType, error) {
	created, err := fs.di.DBRepo().FileRepository.FileType.Create(ctx, data)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}

	return created, nil
}

func (fs *fileType) GetById(ctx context.Context, fileTypeId uint) (*dto.FileType, error) {
	result, err := fs.di.DBRepo().FileRepository.FileType.GetById(ctx, fileTypeId)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}

	return result, nil
}

func (fs *fileType) Update(ctx context.Context, fileTypeId uint, data *dto.FileType) (*dto.FileType, error) {
	updated, err := fs.di.DBRepo().FileRepository.FileType.Update(ctx, fileTypeId, data)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}

	return updated, nil
}

func (fs *fileType) Delete(ctx context.Context, fileTypeId uint) error {
	return fs.di.DBRepo().FileRepository.FileType.Delete(ctx, fileTypeId)
}
