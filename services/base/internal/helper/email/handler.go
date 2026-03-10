package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/sdq-codes/usegro-api/config"
	"github.com/sdq-codes/usegro-api/internal/logger"
)

type MailRequest struct {
	From    string
	To      []string
	Subject string
	Body    string
}

const (
	MIME         = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	devEmailSink = "oyebola.sd@gmail.com"
)

func (r *MailRequest) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.Body = buffer.String()
	return nil
}

func (r *MailRequest) sendMail() (string, error) {
	cfg := config.GetConfig()

	recipients := r.To
	if cfg.Env == "dev" {
		recipients = []string{devEmailSink}
	}

	fromEmail := r.From
	if fromEmail == "" {
		fromEmail = cfg.Ses.FromEmail
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.Ses.Region),
	)
	if err != nil {
		return "", err
	}

	sesClient := ses.NewFromConfig(awsCfg)

	input := &ses.SendEmailInput{
		Source: aws.String(fromEmail),
		Destination: &types.Destination{
			ToAddresses: recipients,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(r.Subject),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data: aws.String(r.Body),
				},
				Text: &types.Content{
					Data: aws.String("Plain text fallback"),
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	out, err := sesClient.SendEmail(ctx, input)
	if err != nil {
		return "", err
	}

	return *out.MessageId, nil
}

func (r *MailRequest) Send(templateName string, items interface{}) {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("The right error %v", err))
	}

	_, err = r.sendMail()
	if err == nil {
		logger.Log.Info("Email successfully sent")
	} else {
		logger.Log.Error(fmt.Sprintf("Error: %v ", err))
	}
}
