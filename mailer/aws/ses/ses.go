package ses

import (
	"fmt"

	"github.com/adeki/go-utils/config"
	"github.com/adeki/go-utils/mailer"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type sesMailer struct {
	client *ses.SES
}

func New() mailer.Mailer {
	c := config.Load()

	sess := session.Must(session.NewSession())
	client := ses.New(sess, aws.NewConfig().WithRegion(c.AWS.Region))

	return &sesMailer{
		client: client,
	}
}

func (m *sesMailer) Send(input mailer.Input) error {
	emailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(input.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(input.Text),
				},
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(input.Html),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(input.Subject),
			},
		},
		Source: aws.String(input.From),
	}

	_, err := m.client.SendEmail(emailInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return fmt.Errorf("%v %w", ses.ErrCodeMessageRejected, aerr)
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return fmt.Errorf("%v %w", ses.ErrCodeMailFromDomainNotVerifiedException, aerr)
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return fmt.Errorf("%v %w", ses.ErrCodeConfigurationSetDoesNotExistException, aerr)
			case ses.ErrCodeConfigurationSetSendingPausedException:
				return fmt.Errorf("%v %w", ses.ErrCodeConfigurationSetSendingPausedException, aerr)
			case ses.ErrCodeAccountSendingPausedException:
				return fmt.Errorf("%v %w", ses.ErrCodeAccountSendingPausedException, aerr)
			default:
				return fmt.Errorf("%v %w", aerr)
			}
		}
		// Return the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return err
	}

	return nil
}
