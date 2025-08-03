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
		logger.Errorf("Failed to init MinIO client: %v", err)
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
		m.logger.Errorf("Failed to check bucket: %v", err)
		return err
	}
	if !exists {
		err = m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			m.logger.Errorf("Failed to create bucket: %v", err)
			return err
		}
	}
	policy := `{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": "*",
			"Action": ["s3:GetObject"],
			"Resource": ["arn:aws:s3:::` + bucket + `/*"]
		}]
	}`
	err = m.client.SetBucketPolicy(ctx, bucket, policy)
	if err != nil {
		m.logger.Errorf("Failed to make public bucket: %v", err)
		return err
	}
	return nil
}

func (m *minioStorage) UploadFile(ctx context.Context, bucket string, objectName string, data io.Reader, size int64) error {
	err := m.NewBucket(ctx, bucket)
	if err != nil {
		return err
	}
	m.logger.Debugf("Put new object %s to bucket %s", objectName, bucket)
	_, err = m.client.PutObject(ctx, bucket, objectName, data, size,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": objectName,
			},
		})
	if err != nil {
		m.logger.Errorf("Failed to upload file: %v", err)
		return err
	}
	return nil
}

func (m *minioStorage) DeleteFile(ctx context.Context, bucket string, objectName string) error {
	err := m.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		m.logger.Errorf("Failed to delete file: %v", err)
		return err
	}
	return nil
}
