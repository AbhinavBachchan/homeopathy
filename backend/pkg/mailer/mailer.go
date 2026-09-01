package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Mailer struct {
	brevoAPIKey string
	senderName  string
	senderEmail string
	httpClient  *http.Client
}

func NewMailer(brevoAPIKey, senderEmail, senderName string) *Mailer {
	if senderEmail == "" {
		senderEmail = "no-reply@homeopathy-platform.com"
	}
	if senderName == "" {
		senderName = "Homeopathy Platform"
	}
	return &Mailer{
		brevoAPIKey: brevoAPIKey,
		senderName:  senderName,
		senderEmail: senderEmail,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

type brevoRecipient struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type brevoSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type brevoEmailPayload struct {
	Sender      brevoSender      `json:"sender"`
	To          []brevoRecipient `json:"to"`
	Subject     string           `json:"subject"`
	HTMLContent string           `json:"htmlContent"`
}

// SendPasswordResetEmail sends the 6-digit password reset code to the user's email.
// If BREVO_API_KEY is not configured (e.g. in local development), it logs the event safely.
func (m *Mailer) SendPasswordResetEmail(recipientEmail, recipientName, resetCode string) error {
	subject := "Your Password Reset Code"
	htmlContent := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 500px; margin: 0 auto; padding: 20px; border: 1px solid #e5e7eb; border-radius: 10px;">
			<h2 style="color: #10847e; margin-top: 0;">Password Reset Request</h2>
			<p>Hello %s,</p>
			<p>We received a request to reset the password for your Homeopathy account. Use the verification code below to complete the reset:</p>
			<div style="background-color: #f0fdf4; border: 1px dashed #10847e; border-radius: 8px; padding: 15px; text-align: center; margin: 20px 0;">
				<span style="font-size: 28px; font-weight: bold; letter-spacing: 6px; color: #10847e;">%s</span>
			</div>
			<p style="color: #4b5563; font-size: 14px;">This code will expire in <strong>15 minutes</strong>.</p>
			<p style="color: #6b7280; font-size: 12px; margin-top: 30px; border-top: 1px solid #e5e7eb; padding-top: 15px;">
				If you did not request a password reset, please ignore this email or contact support if you have concerns.
			</p>
		</div>
	`, recipientName, resetCode)

	if m.brevoAPIKey == "" {
		log.Printf("[MAILER - DEV LOG] Password reset email for %s: Code is [ %s ] (Brevo API key not set)", recipientEmail, resetCode)
		return nil
	}

	payload := brevoEmailPayload{
		Sender: brevoSender{
			Name:  m.senderName,
			Email: m.senderEmail,
		},
		To: []brevoRecipient{
			{
				Email: recipientEmail,
				Name:  recipientName,
			},
		},
		Subject:     subject,
		HTMLContent: htmlContent,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[MAILER ERROR] Failed to encode email payload: %v", err)
		return fmt.Errorf("failed to encode email payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("[MAILER ERROR] Failed to create email request: %v", err)
		return fmt.Errorf("failed to create email request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", m.brevoAPIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		log.Printf("[MAILER ERROR] HTTP dispatch to Brevo failed: %v", err)
		return fmt.Errorf("failed to send email via Brevo: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("[MAILER ERROR] Brevo API rejected email to %s: HTTP %d - Response: %s", recipientEmail, resp.StatusCode, string(respBody))
		return fmt.Errorf("brevo API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[MAILER SUCCESS] Password reset email sent successfully to %s via Brevo", recipientEmail)
	return nil
}
