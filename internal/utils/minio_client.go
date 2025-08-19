package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func uploadFile(ctx context.Context, client *minio.Client, bucketName, localPath, objectName string) error {
	_, err := client.FPutObject(ctx, bucketName, objectName, localPath, minio.PutObjectOptions{})

	if err != nil {
		return fmt.Errorf("failed to upload %s: %w", objectName, err)
	}
	fmt.Println("Uploaded:", objectName)
	return nil
}

func UploadFolderToMinio(localFolderPath, prefix, bucketName, endpoint, accessKey, secretKey string, secure bool, workers int) error {
	ctx := context.Background()

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})

	if err != nil {
		return fmt.Errorf("could not create MinIO client: %w", err)
	}

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("could not check bucket: %w", err)
	}

	if !exists {
		if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("could not create bucket: %w", err)
		}
	}

	var files []struct {
		LocalPath  string
		ObjectName string
	}
	err = filepath.Walk(localFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, _ := filepath.Rel(localFolderPath, path)
			objectName := filepath.ToSlash(filepath.Join(prefix, rel))
			files = append(files, struct {
				LocalPath  string
				ObjectName string
			}{path, objectName})
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not walk folder: %w", err)
	}

	fileChan := make(chan struct {
		LocalPath  string
		ObjectName string
	})

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				if err := uploadFile(ctx, minioClient, bucketName, file.LocalPath, file.ObjectName); err != nil {
					log.Println(err)
				}
			}
		}()
	}

	for _, f := range files {
		fileChan <- f
	}
	close(fileChan)

	wg.Wait()
	masterPath := filepath.Join("movies", localFolderPath, "master.m3u8")
	fmt.Println("Master path:", masterPath)
	return nil
}
