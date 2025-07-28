package minio

import (
	"context"
	"io"

	"github.com/MishaZhem/Gona/src/userservice/internal/app"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

type minioStorage struct {
	client *minio.Client
	logger log.FieldLogger
}

func NewStorageRepository(endpoint, accessKey, secretKey string, useSSL bool, logger log.FieldLogger) app.StorageRepository {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil
	}
	return &minioStorage{
		client: client,
		logger: logger,
	}
}

func (m *minioStorage) NewBucket(ctx context.Context, bucket string) error {
	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		err = m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *minioStorage) UploadFile(bucket string, objectName string, data io.Reader, size int64, contentType string) error {
	return nil
}

func (m *minioStorage) GetFileURL(bucket string, objectName string) (string, error) {
	return "", nil
}

func (m *minioStorage) DeleteFile(bucket string, objectName string) error {
	return nil
}
