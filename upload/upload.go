package upload

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Uploader interface {
	Upload(name string) (string, error)
}

func New(bucket, directory string) Uploader {
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)
	return &uploaderImpl{uploader, bucket, directory}
}

type uploaderImpl struct {
	uploader  *s3manager.Uploader
	bucket    string
	directory string
}

func (u *uploaderImpl) Upload(name string) (string, error) {
	filePath := path.Join(u.directory, name)
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file; %s; %w", filePath, err)
	}
	defer f.Close()
	gracier := "GLACIER"
	result, err := u.uploader.Upload(&s3manager.UploadInput{
		Bucket:       &u.bucket,
		Key:          &name,
		Body:         f,
		StorageClass: &gracier,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object; %w", err)
	}
	return result.Location, nil
}
