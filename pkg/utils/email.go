package utils

import (
	"fmt"
	"net/smtp"
)

// SendEmail sends an HTML formatted email using net/smtp
func SendEmail(smtpHost, smtpPort, smtpUser, smtpPass, sender, recipient, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	header := make(map[string]string)
	header["From"] = sender
	header["To"] = recipient
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(addr, auth, sender, []string{recipient}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	return nil
}
