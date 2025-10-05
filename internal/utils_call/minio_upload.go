package utils_call

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinioClient() (*minio.Client, error) {
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
