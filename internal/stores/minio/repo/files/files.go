package files

import (
	"context"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/bl/use_cases/files/dto"
	"github.com/jkrus/master_api/internal/common/err_const"
	"github.com/jkrus/master_api/pkg/errors"
)

type FileStoreI interface {
	GetFile(ctx context.Context, bucketName, fileUUID string) (*minio.Object, *minio.ObjectInfo, error)
	UploadFile(ctx context.Context, bucketName string, file *dto.FileIN) (string, error)
	DeleteFile(ctx context.Context, bucketName, fileUUID string) error
}

func NewFileStoreI(logger *zap.Logger, minioClient *minio.Client) FileStoreI {
	return &filesRepo{logger: logger, minio: minioClient}
}

type filesRepo struct {
	minio  *minio.Client
	logger *zap.Logger
}

type File struct {
	UUID string
	Name string // Имя файла
	Size int64  // Размер файла

}

func (f *File) fromDTO(dtm *dto.FileIN) *File {
	if dtm == nil {
		return nil
	}

	return &File{
		UUID: dtm.UUID,
		Name: dtm.Name,
		Size: dtm.Size,
	}
}

func (f *File) toDTO() *dto.FileIN {
	if f == nil {
		return nil
	}

	return &dto.FileIN{
		UUID: f.UUID,
		Name: f.Name,
		Size: f.Size,
	}
}

func (f *filesRepo) GetFile(ctx context.Context, bucketName, fileUUID string) (*minio.Object, *minio.ObjectInfo, error) {
	obj, err := f.minio.GetObject(ctx, bucketName, fileUUID, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, errors.Newf("failed to get file with id: %s from minio bucket %s. err: %w", fileUUID, bucketName, err)
	}
	objectInfo, err := obj.Stat()
	if err != nil {
		errMinio := err.(minio.ErrorResponse)
		if errMinio.Message == err_const.ErrKeyDoesNotExist.Error() {
			return nil, nil, errors.Ctx().Just(err_const.ErrMinioRecordNotFound)
		}
		return nil, nil, errors.Ctx().Just(err)
	}

	return obj, &objectInfo, nil
}

func (f *filesRepo) UploadFile(ctx context.Context, bucketName string, file *dto.FileIN) (string, error) {
	exists, errBucketExists := f.minio.BucketExists(ctx, bucketName)
	if errBucketExists != nil || !exists {
		f.logger.Warn(fmt.Sprintf("no bucket %s. creating new one...", bucketName))
		err := f.minio.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", errors.Wrap(err, "failed to create new bucket")
		}
	}

	f.logger.Debug(fmt.Sprintf("put new object %s to bucket %s", file.Name, bucketName))
	_, err := f.minio.PutObject(ctx, bucketName, file.UUID, file.Reader, file.Size,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": file.Name,
			},
			ContentType: "application/octet-stream",
		})
	if err != nil {
		return "", errors.Wrap(err, "failed to upload file")
	}

	return f.generateFileURL(bucketName, file.UUID), nil
}

func (f *filesRepo) DeleteFile(ctx context.Context, bucketName, fileUUID string) error {
	err := f.minio.RemoveObject(ctx, bucketName, fileUUID, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to delete file")
	}
	return nil
}

func (f *filesRepo) generateFileURL(bucketName, fileUUID string) string {
	endpoint := strings.Replace(f.minio.EndpointURL().String(), "localstack", "localhost", -1)
	return fmt.Sprintf("%s/%s/%s", endpoint, bucketName, fileUUID)
}
