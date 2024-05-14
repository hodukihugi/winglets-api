package utils

import (
	"fmt"
	"github.com/hodukihugi/winglets-api/core"
	"net/smtp"
	"sync"
)

func PlainAuth(username, password, host string) smtp.Auth {
	auth := smtp.PlainAuth("GO-AUTH-WINGLETS", username, password, host)
	return auth
}

func SendVerificationEmailAsync(wg *sync.WaitGroup, ch chan error, env *core.Env, verificationToken string, email []string) {
	defer wg.Done()
	a := PlainAuth(env.SmtpUser, env.SmtpPassword, env.SmtpHost)
	subject := fmt.Sprintf("This is your verification token: %s", verificationToken)

	message := []byte("To: " + email[0] + "\r\n" +
		"Subject: Welcome to Winglets website!\r\n" +
		"\r\n" +
		subject + ".\r\n")

	err := smtp.SendMail("smtp.gmail.com:587", a, "Winglets Developer Team", email, message)
	ch <- err
}
