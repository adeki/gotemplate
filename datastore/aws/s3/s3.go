package s3

import (
	"fmt"
	"io"
	"time"

	"github.com/adeki/go-utils/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	svc    s3iface.S3API
	bucket string
}

type Meta struct {
	ContentType string
}

func New(bucket string) *S3 {
	c := config.Load()

	sess := session.Must(session.NewSession())
	svc := s3.New(sess, aws.NewConfig().WithRegion(c.AWS.Region))
	return &S3{
		svc:    svc,
		bucket: bucket,
	}
}

func (mys3 *S3) Put(key string, body io.Reader, contentType string) error {
	uploader := s3manager.NewUploaderWithClient(mys3.svc, func(u *s3manager.Uploader) {
		u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
	})
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(mys3.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	return err
}

func (mys3 *S3) Get(key string) (io.ReadCloser, error) {
	result, err := mys3.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(mys3.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, detectError(err)
	}
	return result.Body, nil
}

func (mys3 *S3) GetPreSignedURL(key string, expire time.Duration) (string, error) {
	req, _ := mys3.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(mys3.bucket),
		Key:    aws.String(key),
	})
	return req.Presign(expire)
}

func (mys3 *S3) Delete(key string) error {
	_, err := mys3.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(mys3.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (mys3 *S3) GetList(prefix string) ([]string, error) {
	result, err := mys3.svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(mys3.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, detectError(err)
	}
	keys := make([]string, len(result.Contents))
	for i, obj := range result.Contents {
		keys[i] = *obj.Key
	}
	return keys, nil
}

func detectError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case s3.ErrCodeNoSuchBucket:
			return fmt.Errorf("%v, %w", s3.ErrCodeNoSuchBucket, aerr)
		case s3.ErrCodeNoSuchKey:
			return fmt.Errorf("%v %w", s3.ErrCodeNoSuchKey, aerr)
		default:
			return aerr
		}
	}
	return err
}
