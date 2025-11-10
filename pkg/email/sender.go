package email

import (
	"email-verify/configs"
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender struct {
	config *configs.Config
}

func NewEmailSender(config *configs.Config) *EmailSender {
	return &EmailSender{
		config: config,
	}
}

func (s *EmailSender) SendVerificationEmail(toEmail, verificationHash string) error {
	// Формируем ссылку для верификации
	verifyURL := fmt.Sprintf("http://localhost:7777/verify/%s", verificationHash)

	// Создаем email
	e := email.NewEmail()
	e.From = s.config.Email.Email
	e.To = []string{toEmail}
	e.Subject = "Подтверждение email"
	e.HTML = []byte(fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Подтверждение Email</title>
		</head>
		<body>
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2 style="color: #333;">Подтвердите ваш email</h2>
				<p>Для завершения верификации вашего email адреса, пожалуйста, перейдите по ссылке ниже:</p>
				<a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #007bff; color: white; text-decoration: none; border-radius: 4px; margin: 16px 0;">
					Подтвердить Email
				</a>
				<p>Или скопируйте эту ссылку в браузер:</p>
				<p style="word-break: break-all; color: #666;">%s</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="color: #999; font-size: 12px;">Если вы не запрашивали подтверждение email, просто проигнорируйте это письмо.</p>
			</div>
		</body>
		</html>
	`, verifyURL, verifyURL))

	// Настройки SMTP для Яндекс
	smtpHost := "smtp.yandex.ru"
	smtpPort := "587" // или 465 для SSL
	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Аутентификация
	auth := smtp.PlainAuth("", s.config.Email.Email, s.config.Email.Password, smtpHost)

	// Отправка письма
	log.Printf("Attempting to send email to: %s", toEmail)
	log.Printf("Using SMTP: %s", smtpAddr)
	log.Printf("From: %s", s.config.Email.Email)

	err := e.Send(smtpAddr, auth)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Verification email sent successfully to: %s", toEmail)
	return nil
}
