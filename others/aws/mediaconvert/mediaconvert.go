package mediaconvert

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/adeki/go-utils/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
)

const (
	// job template names
	mp4ToHlsJobTmpl = "NameOfMp4ToHlsTemplate"
	wavToHlsJobTmpl = "NameOfWavToHlsTemplate"

	// name modifiers
	hlsNameModifier       = "_hls"
	thumbnailNameModifier = "_thumb"
)

func New() (*MediaConvert, error) {
	c := config.Load()

	sess := session.Must(session.NewSession())
	awsConfig := aws.NewConfig().WithRegion(c.AWS.Region)

	tmpsvc := mediaconvert.New(sess, awsConfig)
	tmpmc := &MediaConvert{svc: tmpsvc}
	endpointURLs, err := tmpmc.EndpointURLs()
	if err != nil {
		return nil, err
	}
	svc := mediaconvert.New(sess, awsConfig.WithEndpoint(endpointURLs[0]))
	return &MediaConvert{svc: svc}, nil
}

func makeS3Endpoint(key string) string {
	c := config.Load()
	return "s3://" + c.AWS.S3.Buckets.MediaBucket + "/" + key
}

type MediaConvert struct {
	svc *mediaconvert.MediaConvert
}

func (mc *MediaConvert) EndpointURLs() ([]string, error) {
	result, err := mc.svc.DescribeEndpoints(&mediaconvert.DescribeEndpointsInput{})
	if err != nil {
		return nil, detectError(err)
	}
	urls := make([]string, len(result.Endpoints))
	for i, e := range result.Endpoints {
		urls[i] = *e.Url
	}
	return urls, nil
}

func (mc *MediaConvert) ConvertToHls(inputKey, outputKey string) error {
	switch filepath.Ext(inputKey) {
	case ".wav":
		return mc.WavToHls(inputKey, outputKey)
	case ".mp4":
		return mc.Mp4ToHls(inputKey, outputKey)
	default:
		return errors.New("Invalid file type.")
	}
}

func (mc *MediaConvert) Mp4ToHls(inputKey, outputKey string) error {
	c := config.Load()

	fileInput := makeS3Endpoint(inputKey)
	destination := makeS3Endpoint(outputKey)

	createJobInput := &mediaconvert.CreateJobInput{
		JobTemplate:  aws.String(mp4ToHlsJobTmpl),
		Role:         aws.String(c.AWS.MediaConvert.Role), // maybe no need ?
		UserMetadata: map[string]*string{"job_kind": aws.String("mp4_to_hls")},
		Settings: &mediaconvert.JobSettings{
			Inputs: []*mediaconvert.Input{
				&mediaconvert.Input{
					FileInput: aws.String(fileInput),
				},
			},
			OutputGroups: []*mediaconvert.OutputGroup{
				&mediaconvert.OutputGroup{
					OutputGroupSettings: &mediaconvert.OutputGroupSettings{
						HlsGroupSettings: &mediaconvert.HlsGroupSettings{
							Destination: aws.String(destination),
						},
					},
				},
				&mediaconvert.OutputGroup{
					OutputGroupSettings: &mediaconvert.OutputGroupSettings{
						FileGroupSettings: &mediaconvert.FileGroupSettings{
							Destination: aws.String(destination),
						},
					},
				},
			},
		},
	}

	_, err := mc.svc.CreateJob(createJobInput)
	if err != nil {
		return detectError(err)
	}
	return nil
}

func (mc *MediaConvert) WavToHls(inputKey, outputKey string) error {
	c := config.Load()

	fileInput := makeS3Endpoint(inputKey)
	destination := makeS3Endpoint(outputKey)

	createJobInput := &mediaconvert.CreateJobInput{
		JobTemplate:  aws.String(wavToHlsJobTmpl),
		Role:         aws.String(c.AWS.MediaConvert.Role),
		UserMetadata: map[string]*string{"job_kind": aws.String("wav_to_hls")},
		Settings: &mediaconvert.JobSettings{
			Inputs: []*mediaconvert.Input{
				&mediaconvert.Input{
					FileInput: aws.String(fileInput),
				},
			},
			OutputGroups: []*mediaconvert.OutputGroup{
				&mediaconvert.OutputGroup{
					OutputGroupSettings: &mediaconvert.OutputGroupSettings{
						HlsGroupSettings: &mediaconvert.HlsGroupSettings{
							Destination: aws.String(destination),
						},
					},
				},
			},
		},
	}
	_, err := mc.svc.CreateJob(createJobInput)
	if err != nil {
		return detectError(err)
	}
	return nil
}

func detectError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case mediaconvert.ErrCodeBadRequestException:
			return fmt.Errorf("%v %w", mediaconvert.ErrCodeBadRequestException, aerr)
		case mediaconvert.ErrCodeInternalServerErrorException:
			return fmt.Errorf("%v %w", mediaconvert.ErrCodeInternalServerErrorException, aerr)
		case mediaconvert.ErrCodeForbiddenException:
			return fmt.Errorf("%v %w", mediaconvert.ErrCodeForbiddenException, aerr)
		case mediaconvert.ErrCodeNotFoundException:
			return fmt.Errorf("%v %w", mediaconvert.ErrCodeNotFoundException, aerr)
		case mediaconvert.ErrCodeTooManyRequestsException:
			return fmt.Errorf("%v %w", mediaconvert.ErrCodeTooManyRequestsException, aerr)
		case mediaconvert.ErrCodeConflictException:
			return fmt.Errorf("%v %w", mediaconvert.ErrCodeConflictException, aerr)
		default:
			return aerr
		}
	}
	return err
}
