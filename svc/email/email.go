package email

import (
	"context"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dm"
)

type EmailService struct {
	cli               *dm.Client
	AliyunAccountName string
	FromAlias         string
}

func DefaultEmailService(aliyunAccessKeyId string, aliyunAccessKeySecret string, aliyunAccountName string) *EmailService {
	cli, err := dm.NewClientWithAccessKey("cn-hangzhou", aliyunAccessKeyId, aliyunAccessKeySecret)
	if err != nil {
		panic(err)
	}
	return &EmailService{
		cli:               cli,
		AliyunAccountName: aliyunAccountName,
		FromAlias:         "syzoj-ng",
	}
}

func (s *EmailService) SendEmail(ctx context.Context, To string, Subject string, Body string) error {
	req := dm.CreateSingleSendMailRequest()
	req.Scheme = "https"
	req.AccountName = s.AliyunAccountName
	req.AddressType = requests.NewInteger(1)
	req.ReplyToAddress = requests.NewBoolean(false)
	req.ToAddress = To
	req.FromAlias = s.FromAlias
	req.Subject = Subject
	req.HtmlBody = Body
	_, err := s.cli.SingleSendMail(req)
	if err != nil {
		return err
	}
	return nil
}
