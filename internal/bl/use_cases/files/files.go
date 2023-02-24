package files

import (
	"bytes"
	"context"
	"crypto"
	"io"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/pkg/errors"
)

type FilesI interface {
	Create(ctx context.Context, bucketName string, file *dto.File) (*dto.File, error)
	Delete(ctx context.Context, bucketName, fileUUID string) error
	DownloadFile(ctx context.Context, bucketName, fileUUID string) (f *dto.File, err error)
	Update(ctx context.Context, fileUuid string, data *dto.File) (*dto.File, error)
}

type files struct {
	di internal.IAppDeps
}

func NewFilesI(di internal.IAppDeps) FilesI {
	return &files{di: di}
}

func (f *files) Create(ctx context.Context, bucketName string, file *dto.File) (*dto.File, error) {
	result := &dto.File{}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Ctx().Just(errors.Wrap(err, "failed to generate file uuid"))
	}
	file.Uuid = uuid.String()

	data := make([]byte, file.Size)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}
	_, err = file.Reader.Read(data)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}
	bytes2 := bytes.NewReader(data)
	file.Reader = bytes2
	hash := crypto.SHA256.New()
	hash.Write(data)
	sum := hash.Sum(nil)
	file.CheckSum = sum

	err = f.di.DBRepo().WithTransaction(func(tx *gorm.DB) error {
		result, err = f.di.DBRepo().FileRepository.File.WithTx(tx).Create(ctx, file)
		if err != nil {
			return err
		}

		err = f.di.MinioRepo().Files.FileStore.UploadFile(ctx, bucketName, file)
		if err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}

	return result, nil
}

func (f *files) Delete(ctx context.Context, bucketName, fileUUID string) error {
	err := f.di.DBRepo().WithTransaction(func(tx *gorm.DB) error {
		err := f.di.DBRepo().FileRepository.File.Delete(ctx, fileUUID)
		if err != nil {
			return err
		}

		err = f.di.MinioRepo().Files.FileStore.DeleteFile(ctx, bucketName, fileUUID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Ctx().Just(err)
	}

	return nil
}

func (f *files) DownloadFile(ctx context.Context, bucketName, fileUUID string) (*dto.File, error) {
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
	res := &dto.File{
		Uuid: objectInfo.Key,
		Name: objectInfo.UserMetadata["Name"],
		Data: buffer,
	}
	return res, nil
}

func (f *files) Update(ctx context.Context, fileUuid string, data *dto.File) (*dto.File, error) {
	updated, err := f.di.DBRepo().FileRepository.File.Update(ctx, fileUuid, data)
	if err != nil {
		return nil, errors.Ctx().Just(err)
	}

	return updated, nil
}
