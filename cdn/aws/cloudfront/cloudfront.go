package cloudfront

import (
	"fmt"
	"time"

	"github.com/adeki/go-utils/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

type CloudFront struct {
	svc            *cloudfront.CloudFront
	distributionID string
}

func New() *CloudFront {
	c := config.Load()

	sess := session.Must(session.NewSession())
	svc := cloudfront.New(sess, aws.NewConfig().WithRegion(c.AWS.Region))
	return &CloudFront{
		svc:            svc,
		distributionID: c.AWS.CloudFront.DistributionID,
	}
}

func (cf *CloudFront) Invalidation(paths ...string) error {
	items := make([]*string, len(paths))
	for i, p := range paths {
		items[i] = aws.String(p)
	}
	//https://docs.aws.amazon.com/sdk-for-go/api/service/cloudfront/#CreateInvalidationInput
	input := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(cf.distributionID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(time.Now().Format("20060102150405")),
			Paths: &cloudfront.Paths{
				Items:    items,
				Quantity: aws.Int64(int64(len(items))),
			},
		},
	}
	_, err := cf.svc.CreateInvalidation(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeAccessDenied:
				return fmt.Errorf("%v %w", cloudfront.ErrCodeAccessDenied, aerr)
			case cloudfront.ErrCodeMissingBody:
				return fmt.Errorf("%v %w", cloudfront.ErrCodeMissingBody, aerr)
			case cloudfront.ErrCodeInvalidArgument:
				return fmt.Errorf("%v %w", cloudfront.ErrCodeInvalidArgument, aerr)
			case cloudfront.ErrCodeNoSuchDistribution:
				return fmt.Errorf("%v %w", cloudfront.ErrCodeNoSuchDistribution, aerr)
			case cloudfront.ErrCodeBatchTooLarge:
				return fmt.Errorf("%v %w", cloudfront.ErrCodeBatchTooLarge, aerr)
			case cloudfront.ErrCodeTooManyInvalidationsInProgress:
				return fmt.Errorf("%v %w", cloudfront.ErrCodeTooManyInvalidationsInProgress, aerr)
			case cloudfront.ErrCodeInconsistentQuantities:
				return fmt.Errorf("%v %w", cloudfront.ErrCodeInconsistentQuantities, aerr)
			default:
				return aerr
			}
		}
		// Return the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return err
	}
	return nil
}
