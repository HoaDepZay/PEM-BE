package minio

import (
	"context"
	"io"
	"log"
	"os"

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

	// Make a new bucket called mymusic.
	ctx := context.Background()
	err = MinioClient.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := MinioClient.BucketExists(ctx, BucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", BucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", BucketName)
		
		// Set public read policy for the bucket
		policy := `{"Version":"2012-10-17","Statement":[{"Action":["s3:GetObject"],"Effect":"Allow","Principal":{"AWS":["*"]},"Resource":["arn:aws:s3:::` + BucketName + `/*"]}]}`
		err = MinioClient.SetBucketPolicy(ctx, BucketName, policy)
		if err != nil {
			log.Println("Failed to set bucket policy:", err)
		}
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
