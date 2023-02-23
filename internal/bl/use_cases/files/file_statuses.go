package files

import (
	"context"

	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/pkg/errors"
)

type FileStatusesI interface {
	Create(ctx context.Context, data *dto.FileStatus) (*dto.FileStatus, error)
	GetById(ctx context.Context, fileStatusId uint) (*dto.FileStatus, error)
	Update(ctx context.Context, fileStatusId uint, data *dto.FileStatus) (*dto.FileStatus, error)
	Delete(ctx context.Context, fileStatusId uint) error
}

type fileStatus struct {
	di internal.IAppDeps
}

func NewFileStatusI(di internal.IAppDeps) FileStatusesI {
	return &fileStatus{
		di: di,
	}
}

func (fs *fileStatus) Create(ctx context.Context, data *dto.FileStatus) (*dto.FileStatus, error) {
	created, err := fs.di.DBRepo().FileRepository.FileStatusRepository.Create(ctx, data)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}

	return created, nil
}

func (fs *fileStatus) GetById(ctx context.Context, fileStatusId uint) (*dto.FileStatus, error) {
	result, err := fs.di.DBRepo().FileRepository.FileStatusRepository.GetById(ctx, fileStatusId)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}

	return result, nil
}

func (fs *fileStatus) Update(ctx context.Context, fileStatusId uint, data *dto.FileStatus) (*dto.FileStatus, error) {
	updated, err := fs.di.DBRepo().FileRepository.FileStatusRepository.Update(ctx, fileStatusId, data)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}

	return updated, nil
}

func (fs *fileStatus) Delete(ctx context.Context, fileStatusId uint) error {
	return fs.di.DBRepo().FileRepository.FileStatusRepository.Delete(ctx, fileStatusId)
}
