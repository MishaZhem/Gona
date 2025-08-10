package minio

import (
	"context"
	"io"
	"strings"

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
	endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "https://"), "http://")

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
		m.logger.Infof("Bucket '%s' does not exist. Creating...", bucket)
		err = m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			m.logger.Errorf("Failed to create bucket: %v", err)
			return err
		}
	}

	m.logger.Infof("Try to make public bucket")
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
	m.logger.Infof("Public policy applied to bucket '%s'", bucket)
	return nil
}

func (m *minioStorage) UploadFile(ctx context.Context, bucket string, objectName string, data io.Reader, size int64, contentType string) error {
	err := m.NewBucket(ctx, bucket)
	if err != nil {
		return err
	}
	m.logger.Infof("Put new object %s to bucket %s", objectName, bucket)
	m.logger.Infof("Upload file: objectName=%s, fileSize=%d", objectName, size)
	_, err = m.client.PutObject(ctx, bucket, objectName, data, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		m.logger.Errorf("Failed to upload file: %v", err)
		return err
	}

	stat, err := m.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		m.logger.Errorf("stat after upload failed: %w", err)
		return err
	}
	m.logger.Infof("Stat OK: key=%s size=%d contentType=%s", objectName, stat.Size, stat.ContentType)
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
