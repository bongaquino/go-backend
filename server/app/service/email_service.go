package service

import "koneksi/server/app/provider"

type EmailService struct {
	postmarkProvider *provider.PostmarkProvider
}

func NewEmailService(postmarkProvider *provider.PostmarkProvider) *EmailService {
	return &EmailService{
		postmarkProvider: postmarkProvider,
	}
}

func (es *EmailService) SendWelcomeEmail(to string) error {
	subject := "Welcome to Koneksi!"
	body := "<h1>Welcome to Koneksi</h1><p>Thank you for joining us!</p>"
	return es.postmarkProvider.SendEmail(to, subject, body)
}

func (es *EmailService) SendVerificationCode(to, code string) error {
	subject := "Verification Code"
	body := "<h1>Verification Code</h1><p>Your verification code is: " + code + "</p>"
	return es.postmarkProvider.SendEmail(to, subject, body)
}

func (es *EmailService) SendPasswordResetCode(to, code string) error {
	subject := "Password Reset Code"
	body := "<h1>Password Reset Code</h1><p>Your password reset code is: " + code + "</p>"
	return es.postmarkProvider.SendEmail(to, subject, body)
}
func (es *EmailService) SendFileShareNotification(to, fileName string) error {
	subject := "File Shared: " + fileName
	body := "<h1>File Shared</h1><p>A file named '" + fileName + "' has been shared with you.</p>"
	return es.postmarkProvider.SendEmail(to, subject, body)
}
