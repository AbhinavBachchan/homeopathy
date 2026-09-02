package config

import "os"

// Config holds all environment-driven settings.
// Names mirror the "Environment Variables Required" section of the dev brief
// so nothing needs to be renamed when wiring up Razorpay/Stripe/Interakt/etc.
type Config struct {
	Port string
	Env  string

	DatabaseURL string

	JWTSecret       string
	JWTExpiryHours  int

	RazorpayKeyID     string
	RazorpayKeySecret string

	StripePublishableKey string
	StripeSecretKey      string
	StripeWebhookSecret  string

	AnthropicAPIKey string // Claude API - symptom assessment, chatbot, content
	OpenAIAPIKey    string // GPT-4o fallback

	InteraktAPIKey string // WhatsApp

	BrevoAPIKey      string // Transactional + marketing email
	BrevoSenderEmail string // Verified sender email in Brevo
	BrevoSenderName  string // Sender display name

	MSG91AuthKey     string // OTP
	MSG91TemplateID  string

	ShiprocketEmail    string
	ShiprocketPassword string

	AlgoliaAppID    string
	AlgoliaAdminKey string

	N8NWebhookSecret string

	SentryDSN string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("APP_ENV", "development"),

		DatabaseURL: getEnv("DATABASE_URL", "postgresql://postgres:Abhi9065890801@db.fsncoziyuqtwighkxait.supabase.co:5432/postgres"),

		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),

		RazorpayKeyID:     getEnv("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret: getEnv("RAZORPAY_KEY_SECRET", ""),

		StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),

		AnthropicAPIKey: getEnv("ANTHROPIC_API_KEY", ""),
		OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),

		InteraktAPIKey: getEnv("INTERAKT_API_KEY", ""),

		BrevoAPIKey:      getEnv("BREVO_API_KEY", ""),
		BrevoSenderEmail: getEnv("BREVO_SENDER_EMAIL", ""),
		BrevoSenderName:  getEnv("BREVO_SENDER_NAME", "Homeopathy Platform"),

		MSG91AuthKey:    getEnv("MSG91_AUTH_KEY", ""),
		MSG91TemplateID: getEnv("MSG91_TEMPLATE_ID", ""),

		ShiprocketEmail:    getEnv("SHIPROCKET_EMAIL", ""),
		ShiprocketPassword: getEnv("SHIPROCKET_PASSWORD", ""),

		AlgoliaAppID:    getEnv("ALGOLIA_APP_ID", ""),
		AlgoliaAdminKey: getEnv("ALGOLIA_ADMIN_KEY", ""),

		N8NWebhookSecret: getEnv("N8N_WEBHOOK_SECRET", ""),

		SentryDSN: getEnv("SENTRY_DSN", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
