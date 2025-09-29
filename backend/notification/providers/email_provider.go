package providers

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"tchat.dev/notification/models"
)

// EmailProvider implements ChannelProvider for email notifications
type EmailProvider struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
}

// EmailConfig holds email provider configuration
type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     string `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	FromEmail    string `json:"from_email"`
	FromName     string `json:"from_name"`
}

// NewEmailProvider creates a new email provider
func NewEmailProvider(config EmailConfig) *EmailProvider {
	return &EmailProvider{
		smtpHost:     config.SMTPHost,
		smtpPort:     config.SMTPPort,
		smtpUsername: config.SMTPUsername,
		smtpPassword: config.SMTPPassword,
		fromEmail:    config.FromEmail,
		fromName:     config.FromName,
	}
}

// Send sends an email notification
func (e *EmailProvider) Send(ctx context.Context, notification *models.Notification) error {
	if notification.RecipientEmail == "" {
		return fmt.Errorf("recipient email is required for email notifications")
	}

	// Create email content
	subject := notification.Title
	if subject == "" {
		subject = "Notification"
	}

	body := e.buildEmailBody(notification)

	// Create message
	message := e.buildMessage(e.fromEmail, notification.RecipientEmail, subject, body)

	// Send email
	auth := smtp.PlainAuth("", e.smtpUsername, e.smtpPassword, e.smtpHost)
	addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)

	err := smtp.SendMail(addr, auth, e.fromEmail, []string{notification.RecipientEmail}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendBatch sends multiple email notifications
func (e *EmailProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	for _, notification := range notifications {
		if err := e.Send(ctx, notification); err != nil {
			// Log error but continue with other notifications
			// In a real implementation, you might want to collect errors and return them
			continue
		}
	}
	return nil
}

// ValidateConfig validates the email provider configuration
func (e *EmailProvider) ValidateConfig() error {
	if e.smtpHost == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if e.smtpPort == "" {
		return fmt.Errorf("SMTP port is required")
	}
	if e.smtpUsername == "" {
		return fmt.Errorf("SMTP username is required")
	}
	if e.smtpPassword == "" {
		return fmt.Errorf("SMTP password is required")
	}
	if e.fromEmail == "" {
		return fmt.Errorf("from email is required")
	}
	return nil
}

// GetProviderName returns the provider name
func (e *EmailProvider) GetProviderName() string {
	return "email"
}

// SupportsChannels returns the channels this provider supports
func (e *EmailProvider) SupportsChannels() []models.NotificationChannel {
	return []models.NotificationChannel{models.NotificationChannelEmail}
}

// buildEmailBody builds the email body content
func (e *EmailProvider) buildEmailBody(notification *models.Notification) string {
	var body strings.Builder

	// Add content
	if notification.Body != "" {
		body.WriteString(notification.Body)
	}

	// Add media if present
	if notification.MediaURL != "" {
		body.WriteString(fmt.Sprintf("\n\nAttachment: %s", notification.MediaURL))
	}

	// Add action buttons if present
	if notification.ActionURL != "" {
		if notification.ActionText != "" {
			body.WriteString(fmt.Sprintf("\n\n[%s](%s)", notification.ActionText, notification.ActionURL))
		} else {
			body.WriteString(fmt.Sprintf("\n\nLink: %s", notification.ActionURL))
		}
	}

	return body.String()
}

// buildMessage builds the email message
func (e *EmailProvider) buildMessage(from, to, subject, body string) string {
	var message strings.Builder

	// Headers
	message.WriteString(fmt.Sprintf("From: %s <%s>\n", e.fromName, from))
	message.WriteString(fmt.Sprintf("To: %s\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	message.WriteString("MIME-Version: 1.0\n")
	message.WriteString("Content-Type: text/plain; charset=utf-8\n")
	message.WriteString("\n")

	// Body
	message.WriteString(body)

	return message.String()
}

// SendGridEmailProvider implements email provider using SendGrid
type SendGridEmailProvider struct {
	apiKey   string
	fromEmail string
	fromName  string
}

// NewSendGridEmailProvider creates a new SendGrid email provider
func NewSendGridEmailProvider(apiKey, fromEmail, fromName string) *SendGridEmailProvider {
	return &SendGridEmailProvider{
		apiKey:   apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

// Send sends an email via SendGrid
func (s *SendGridEmailProvider) Send(ctx context.Context, notification *models.Notification) error {
	// This would use the SendGrid API to send emails
	// For now, it's a placeholder that logs the action
	fmt.Printf("SendGrid: Sending email to %s with subject '%s'\n",
		notification.RecipientEmail, notification.Title)
	return nil
}

// SendBatch sends multiple emails via SendGrid
func (s *SendGridEmailProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	for _, notification := range notifications {
		if err := s.Send(ctx, notification); err != nil {
			continue
		}
	}
	return nil
}

// ValidateConfig validates SendGrid configuration
func (s *SendGridEmailProvider) ValidateConfig() error {
	if s.apiKey == "" {
		return fmt.Errorf("SendGrid API key is required")
	}
	if s.fromEmail == "" {
		return fmt.Errorf("from email is required")
	}
	return nil
}

// GetProviderName returns the provider name
func (s *SendGridEmailProvider) GetProviderName() string {
	return "sendgrid"
}

// SupportsChannels returns supported channels
func (s *SendGridEmailProvider) SupportsChannels() []models.NotificationChannel {
	return []models.NotificationChannel{models.NotificationChannelEmail}
}