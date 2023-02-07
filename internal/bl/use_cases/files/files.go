package files

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/pkg/errors"
)

type FilesI interface {
	Create(ctx context.Context, bucketName string, file *dto.FileIN) (string, error)
	Delete(ctx context.Context, bucketName, fileUUID string) error
	GetFile(ctx context.Context, bucketName, fileUUID string) (f *dto.FileOUT, err error)
}

type files struct {
	di internal.IAppDeps
}

func NewFilesI(di internal.IAppDeps) FilesI {
	return &files{di: di}
}

func (f *files) Create(ctx context.Context, bucketName string, file *dto.FileIN) (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", errors.Ctx().Just(errors.Wrap(err, "failed to generate file uuid"))
	}
	file.UUID = uuid.String()

	url, err := f.di.MinioRepo().Files.FileStore.UploadFile(ctx, bucketName, file)
	if err != nil {
		return "", errors.Ctx().Just(err)
	}
	return url, nil
}

func (f *files) Delete(ctx context.Context, bucketName, fileUUID string) error {
	return f.di.MinioRepo().Files.FileStore.DeleteFile(ctx, bucketName, fileUUID)
}

func (f *files) GetFile(ctx context.Context, bucketName, fileUUID string) (*dto.FileOUT, error) {
	obj, objectInfo, err := f.di.MinioRepo().Files.FileStore.GetFile(ctx, bucketName, fileUUID)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}
	defer obj.Close()
	buffer := make([]byte, objectInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, errors.Ctx().Just(err)
	}
	res := &dto.FileOUT{
		UUID:  objectInfo.Key,
		Name:  objectInfo.UserMetadata["Name"],
		Size:  objectInfo.Size,
		Bytes: buffer,
	}
	return res, nil
}
