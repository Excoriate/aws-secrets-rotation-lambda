package s3

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/adapter"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/common"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/daggerio"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/erroer"
	"golang.org/x/exp/rand"
	"os"
)

type Manager interface {
	UploadFile(bucket, key, filePath string) error
}

type Bucket struct {
	Client     *dagger.Client
	Ctx        context.Context
	S3         *s3.Client
	S3Uploader *manager.Uploader
}

func GenerateRandomObjectSuffix() string {
	rand.Seed(uint64(rand.Int()))
	return fmt.Sprintf("%d", rand.Int())
}

func (b *Bucket) UploadFile(bucket, key, filePath string) error {
	_, err := b.S3.HeadObject(b.Ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	// No error means the destination file already exist.
	if err == nil {

		//// Create a new destination path with a random suffix.
		//newKey := fmt.Sprintf("%s-%s", key, GenerateRandomObjectSuffix())
		//
		//_, err := client.CopyObject(b.Ctx, &s3.CopyObjectInput{
		//	Bucket:     &bucket,
		//	CopySource: &key,
		//	Key:        &newKey,
		//})

		//if err != nil {
		//	return erroer.NewPipelineConfigurationError(
		//		fmt.Sprintf("Failed to replace existing file %s with the new generated prefix %s", key, newKey),
		//		err)
		//}

		_, err = b.S3.DeleteObject(b.Ctx, &s3.DeleteObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})

		if err != nil {
			return erroer.NewPipelineConfigurationError(
				fmt.Sprintf("Failed to delete existing file %s", key),
				err)
		}
	}

	err = common.FileExist(filePath)
	if err != nil {
		return erroer.NewPipelineConfigurationError(
			fmt.Sprintf("File %s does not exist", filePath),
			err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return erroer.NewPipelineConfigurationError(
			fmt.Sprintf("Failed to open file %s", filePath),
			err)
	}

	defer file.Close()

	_, err = b.S3Uploader.Upload(b.Ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	})

	if err != nil {
		return erroer.NewPipelineConfigurationError(
			fmt.Sprintf("Failed to upload file %s to bucket %s with key %s", filePath, bucket, key),
			err)
	}

	return nil
}

func NewS3(client *dagger.Client, ctx context.Context) (*Bucket, error) {
	awsCreds := daggerio.ScanAWSCredsFromEnv()

	if len(awsCreds) == 0 {
		return nil, erroer.NewPipelineConfigurationError(
			"AWS credentials not found in environment",
			nil)
	}

	cfg, err := adapter.NewAWS("us-east-1") // FIXME: Make this configurable.
	if err != nil {
		return nil, erroer.NewPipelineConfigurationError(
			"Failed to create AWS adapter",
			err)
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)

	return &Bucket{
		Client:     client,
		Ctx:        ctx,
		S3:         s3Client,
		S3Uploader: uploader,
	}, nil
}
