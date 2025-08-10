package provider

import (
	"log"

	"bongaquino/server/config"

	"github.com/keighl/postmark"
)

// PostmarkProvider handles email sending using Postmark
type PostmarkProvider struct {
	client *postmark.Client
}

// NewPostmarkProvider initializes a new PostmarkProvider
func NewPostmarkProvider() *PostmarkProvider {
	postmarkConfig := config.LoadPostmarkConfig()

	client := postmark.NewClient(postmarkConfig.PostmarkAPIKey, "")
	return &PostmarkProvider{
		client: client,
	}
}

// SendEmail sends an email using Postmark
func (p *PostmarkProvider) SendEmail(to, subject, body string) error {
	postmarkConfig := config.LoadPostmarkConfig()

	email := postmark.Email{
		From:     postmarkConfig.PostmarkFrom,
		To:       to,
		Subject:  subject,
		HtmlBody: body,
	}

	_, err := p.client.SendEmail(email)
	if err != nil {
		log.Printf("failed to send email: %v", err)
		return err
	}

	return nil
}
