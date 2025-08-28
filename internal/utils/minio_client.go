package utils

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClientPool struct {
	clients map[string]*minio.Client
	mutex   sync.RWMutex
}

var (
	clientPool = &MinIOClientPool{
		clients: make(map[string]*minio.Client),
	}
)

func createHTTPTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    false,
	}
}

func GetMinIOClient(endpoint, accessKey, secretKey string, secure bool) (*minio.Client, error) {
	clientKey := fmt.Sprintf("%s:%s:%v", endpoint, accessKey, secure)

	clientPool.mutex.RLock()
	client, exists := clientPool.clients[clientKey]
	clientPool.mutex.RUnlock()

	if exists {
		return client, nil
	}

	clientPool.mutex.Lock()
	defer clientPool.mutex.Unlock()

	if client, exists := clientPool.clients[clientKey]; exists {
		return client, nil
	}

	transport := createHTTPTransport()

	client, err := minio.New(endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:    secure,
		Transport: transport,
	})

	if err != nil {
		return nil, fmt.Errorf("could not create MinIO client: %w", err)
	}

	clientPool.clients[clientKey] = client
	return client, nil
}

func CleanupMinIOClientPool() {
	clientPool.mutex.Lock()
	defer clientPool.mutex.Unlock()

	for key := range clientPool.clients {
		delete(clientPool.clients, key)
	}

	fmt.Println("🧹 MinIO client pool cleaned up")
}

func retryWithBackoff(ctx context.Context, maxRetries int, operation func() error) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {

			delay := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}

			fmt.Printf("Retrying in %v (attempt %d/%d)...\n", delay, attempt, maxRetries)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		if ctx.Err() != nil {
			return ctx.Err()
		}

		fmt.Printf("Upload attempt %d failed: %v\n", attempt+1, err)
	}

	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries+1, lastErr)
}

func uploadFile(ctx context.Context, client *minio.Client, bucketName, localPath, objectName string) error {
	operation := func() error {
		_, err := client.FPutObject(ctx, bucketName, objectName, localPath, minio.PutObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to upload %s: %w", objectName, err)
		}
		return nil
	}

	err := retryWithBackoff(ctx, 3, operation)
	if err != nil {
		return err
	}

	fmt.Println("✅ Uploaded:", objectName)
	return nil
}

func UploadFolderToMinio(localFolderPath, prefix, bucketName, endpoint, accessKey, secretKey string, secure bool, workers int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	minioClient, err := GetMinIOClient(endpoint, accessKey, secretKey, secure)
	if err != nil {
		return fmt.Errorf("could not get MinIO client: %w", err)
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

	errorChan := make(chan error, len(files))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for file := range fileChan {
				fmt.Printf("🚀 Worker %d uploading: %s\n", workerID, file.ObjectName)
				if err := uploadFile(ctx, minioClient, bucketName, file.LocalPath, file.ObjectName); err != nil {
					log.Printf("❌ Worker %d failed to upload %s: %v", workerID, file.ObjectName, err)
					select {
					case errorChan <- fmt.Errorf("failed to upload %s: %w", file.ObjectName, err):
					default:

					}
				}
			}
		}(i)
	}

	for _, f := range files {
		fileChan <- f
	}
	close(fileChan)

	wg.Wait()
	close(errorChan)

	var uploadErrors []error
	for err := range errorChan {
		uploadErrors = append(uploadErrors, err)
	}

	if len(uploadErrors) > 0 {
		fmt.Printf("⚠️  Upload completed with %d errors:\n", len(uploadErrors))
		for _, err := range uploadErrors {
			fmt.Printf("   • %v\n", err)
		}

		return fmt.Errorf("upload failed with %d errors, first error: %w", len(uploadErrors), uploadErrors[0])
	}

	return nil
}
