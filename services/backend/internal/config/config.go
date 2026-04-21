package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

const (
	ENV_DEVELOPMENT = "development"
	ENV_PRODUCTION  = "production"
)

// Config holds all configuration for the application
type Config struct {
	Server             ServerConfig
	Auth               AuthConfig
	ImageConf          S3Config `env:",prefix=IMAGE_"`
	ImageDomain        string   `env:"IMAGE_DOMAIN"`
	MongodbURL         string   `env:"MONGODB_URL"`
	MongodbDatabase    string   `env:"MONGODB_DATABASE"`
	LL2URLPrefix       string   `env:"LL2_URL_PREFIX"`
	LL2RequestInterval int      `env:"LL2_REQUEST_INTERVAL, default=5"` // in seconds
	Email              EmailConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `env:"SERVER_PORT"`
	Host string `env:"SERVER_HOST"`
	Env  string `env:"SERVER_ENVIRONMENT"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret              string `env:"JWT_SECRET,required"`
	AccessTokenExpireMin   int    `env:"ACCESS_TOKEN_EXPIRE_MIN,default=15"`  // in minutes
	RefreshTokenExpireDays int    `env:"REFRESH_TOKEN_EXPIRE_DAYS,default=7"` // in days
	Issuer                 string `env:"JWT_ISSUER,default=launchdate-backend"`
}

// S3Config holds S3 configuration
type S3Config struct {
	Endpoint  string `env:"S3_ENDPOINT"`
	Region    string `env:"S3_REGION"`
	Bucket    string `env:"S3_BUCKET"`
	AccessKey string `env:"S3_ACCESS_KEY"`
	SecretKey string `env:"S3_SECRET_KEY"`
	Domain    string `env:"S3_DOMAIN"`
}

// EmailConfig holds email configuration
type EmailConfig struct {
	ResendAPIKey string `env:"RESEND_API_KEY"`
	From         string `env:"EMAIL_FROM"`
	WebBaseURL   string `env:"EMAIL_WEB_BASE_URL"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{}
	if err := envconfig.Process(context.Background(), config); err != nil {
		return nil, err
	}

	return config, nil
}
