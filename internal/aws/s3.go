package aws

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Store struct {
	Bucket  string
	Session *session.Session
}

func NewS3Store(bucket string) (*S3Store, error) {
	sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})
	if err != nil {
		return nil, fmt.Errorf("can't create aws session (%s)", err)
	}

	return &S3Store{Bucket: bucket, Session: sess}, nil
}

func (s *S3Store) Download(key string) ([]byte, error) {
	s3Svc := s3.New(s.Session)
	result, err := s3Svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("can't get s3 headobject (%s)", err)
	}

	buff := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(s.Session)
	n, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("can't download object from s3 (s3://%s/%s) (%s)", s.Bucket, key, err)
	}

	if n != aws.Int64Value(result.ContentLength) {
		return nil, fmt.Errorf("can't download whole object from s3 (s3://%s/%s) size (%d)", s.Bucket, key, n)
	}

	if n < 1 {
		return nil, fmt.Errorf("zero bytes written to memory")
	}

	return buff.Bytes(), nil
}

func (s *S3Store) Upload(key string, b []byte) error {
	uploader := s3manager.NewUploader(s.Session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(b),
	})
	if err != nil {
		return fmt.Errorf("can't upload object to s3 (s3://%s/%s) (%s)", s.Bucket, key, err)
	}

	return nil
}
