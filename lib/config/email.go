package config

import (
	"fmt"
	"os"

	"github.com/syzoj/syzoj-ng-go/svc/email"
)

func NewAliyunEmail(prefix string) (*email.EmailService, error) {
	accessKey := os.Getenv(prefix + "ACCESS_KEY_ID")
	secret := os.Getenv(prefix + "ACCESS_KEY_SECRET")
	accountName := os.Getenv(prefix + "ACCOUNT_NAME")
	if accessKey == "" || secret == "" || accountName == "" {
		return nil, fmt.Errorf("%sACCESS_KEY_ID or %sACCESS_KEY_SECRET or %sACCOUNT_NAME missing", prefix, prefix, prefix)
	}
	return email.DefaultEmailService(accessKey, secret, accountName), nil
}
