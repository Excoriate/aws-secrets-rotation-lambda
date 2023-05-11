package tasks

import (
	"context"
	"fmt"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/config"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/daggerio"
	s32 "github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/s3"
	"github.com/excoriate/aws-secrets-rotation-lambda/dagger-pipeline/internal/tui"
	"strings"
)

func UploadToS3() error {
	ux := tui.NewTitle()
	msg := tui.NewTUIMessage()
	ux.ShowSubTitle("lambda:", "PushS3")

	// Getting working directories.
	dirs, err := config.GetDirConfig()
	if err != nil {
		msg.ShowError("", "Failed to get working directories", err)
		return err
	}

	msg.ShowInfo("", fmt.Sprintf("Current root directory: %s", dirs.CurrentDir))

	// Booting dagger!
	ctx := context.Background()
	client, err := daggerio.NewClient(ctx)
	if err != nil {
		msg.ShowError("", "Failed to boot dagger", err)
		return err
	}

	// Fetching configuration from Viper.
	cfg := config.Cfg{}
	s3BucketCfg, s3BucketCfgErr := cfg.GetFromViper("s3-bucket")
	if s3BucketCfgErr != nil {
		msg.ShowError("", "Failed to get s3 bucket configuration", s3BucketCfgErr)
		return s3BucketCfgErr
	}

	s3Bucket := s3BucketCfg.Value.(string)
	s3DestinationPathCfg, s3DestinationPathErr := cfg.GetFromViper("s3-destination-path")
	if s3DestinationPathErr != nil {
		msg.ShowError("", "Failed to get s3 destination path", s3DestinationPathErr)
		return s3DestinationPathErr
	}

	s3DestinationPath := s3DestinationPathCfg.Value.(string)

	s3FileToUploadCfg, s3FileToUploadCfgErr := cfg.GetFromViper("s3-file-to-upload")
	if s3FileToUploadCfgErr != nil {
		msg.ShowError("", "Failed to get s3 file to upload", s3FileToUploadCfgErr)
		return s3FileToUploadCfgErr
	}

	s3FileToUpload := s3FileToUploadCfg.Value.(string)

	msg.ShowInfo("", fmt.Sprintf("Task S3-Upload will upload the file %s to the bucket %s in the"+
		" destination path %s", s3FileToUpload, s3Bucket, s3DestinationPath))

	// Check if the destination path doesn't includes a file. If it doesn't,
	//add the filename of the s3FileToUpload to the destination path.
	if !strings.Contains(s3DestinationPath, "/") {
		msg.ShowWarning("", "The destination path doesn't include a file. "+
			"Adding the filename of the"+" s3FileToUpload to the destination path.")

		// The s3FileToUpload is a path. The filename always is in the last part of the path.
		//Add it smartly.
		s3DestinationPath = fmt.Sprintf("%s/%s", s3DestinationPath, strings.Split(s3FileToUpload, "/")[len(strings.Split(s3FileToUpload, "/"))-1])

		msg.ShowInfo("", fmt.Sprintf("The new destination path is %s", s3DestinationPath))
	}

	s3, uploadErr := s32.NewS3(client, ctx)
	if uploadErr != nil {
		msg.ShowError("", "Failed to create S3 client", uploadErr)
		return uploadErr
	}

	if s3.UploadFile(s3Bucket, s3DestinationPath, s3FileToUpload) != nil {
		msg.ShowError("", "Failed to upload file to S3", err)
		return err
	}

	msg.ShowSuccess("", "Successfully uploaded file to S3")

	defer client.Close()

	return nil
}
