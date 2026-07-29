package minio

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client
var BucketName string

func ConnectMinIO() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "minio:9000" // Default for docker-compose internal network
	}
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	if accessKeyID == "" {
		accessKeyID = "admin"
	}
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	if secretAccessKey == "" {
		secretAccessKey = "password123"
	}
	BucketName = os.Getenv("MINIO_BUCKET_NAME")
	if BucketName == "" {
		BucketName = "visual-finance-bucket"
	}
	useSSL := false // Set to true if MinIO is configured with HTTPS

	// Initialize minio client object.
	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("Failed to connect to MinIO:", err)
	}

	log.Printf("Successfully connected to MinIO %s\n", endpoint)

	// Make a new bucket called visual-finance-bucket.
	ctx := context.Background()
	var err error

	// Retry loop for MakeBucket (wait for MinIO to start)
	for i := 0; i < 5; i++ {
		err = MinioClient.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
		if err == nil {
			break
		}
		// Check to see if we already own this bucket
		exists, errBucketExists := MinioClient.BucketExists(ctx, BucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", BucketName)
			err = nil
			break
		}
		log.Printf("Waiting for MinIO to initialize (retry %d/5)...\n", i+1)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Println("Failed to create bucket after retries, skipping policy setup:", err)
		return
	}

	log.Printf("Successfully created %s\n", BucketName)
	
	// Set public read policy for the bucket
	policy := `{"Version":"2012-10-17","Statement":[{"Action":["s3:GetObject"],"Effect":"Allow","Principal":{"AWS":["*"]},"Resource":["arn:aws:s3:::` + BucketName + `/*"]}]}`
	err = MinioClient.SetBucketPolicy(ctx, BucketName, policy)
	if err != nil {
		log.Println("Failed to set bucket policy:", err)
	}
}

// UploadImage uploads an image to MinIO and returns the URL
func UploadImage(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	_, err := MinioClient.PutObject(ctx, BucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	
	// Trả về relative URL hoặc public URL của MinIO (tuỳ cách cấu hình reverse proxy)
	return "/" + BucketName + "/" + objectName, nil
}
