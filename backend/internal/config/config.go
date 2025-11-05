package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerHost string
	ServerPort string
	Environment string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSL      string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Qdrant
	QdrantHost string
	QdrantPort string
	QdrantAPIKey string

	// JWT
	JWTSecret              string
	JWTAccessTokenMinutes  int
	JWTRefreshTokenDays    int

	// AWS/Storage
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	S3BucketName       string

	// Email
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// WebRTC
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioAPIKey     string
	TwilioAPISecret  string

	// OpenAI
	OpenAIAPIKey string
	OpenAIModel  string

	// App Settings
	MaxUploadSize    int64
	AllowedOrigins   []string
	RateLimitPerMin  int
	ProfileVideoMaxDuration int
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	cfg := &Config{
		// Server
		ServerHost:  getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "findme"),
		DBPassword: getEnv("DB_PASSWORD", "findme_password"),
		DBName:     getEnv("DB_NAME", "findme_db"),
		DBSSL:      getEnv("DB_SSL", "disable"),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		// Qdrant
		QdrantHost:   getEnv("QDRANT_HOST", "localhost"),
		QdrantPort:   getEnv("QDRANT_PORT", "6333"),
		QdrantAPIKey: getEnv("QDRANT_API_KEY", ""),

		// JWT
		JWTSecret:             getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTAccessTokenMinutes: getEnvInt("JWT_ACCESS_TOKEN_MINUTES", 15),
		JWTRefreshTokenDays:   getEnvInt("JWT_REFRESH_TOKEN_DAYS", 7),

		// AWS
		AWSRegion:          getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		S3BucketName:       getEnv("S3_BUCKET_NAME", ""),

		// Email
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@findme.app"),

		// WebRTC
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioAPIKey:     getEnv("TWILIO_API_KEY", ""),
		TwilioAPISecret:  getEnv("TWILIO_API_SECRET", ""),

		// OpenAI
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:  getEnv("OPENAI_MODEL", "gpt-4"),

		// App Settings
		MaxUploadSize:           getEnvInt64("MAX_UPLOAD_SIZE", 100*1024*1024), // 100MB
		RateLimitPerMin:         getEnvInt("RATE_LIMIT_PER_MIN", 60),
		ProfileVideoMaxDuration: getEnvInt("PROFILE_VIDEO_MAX_DURATION", 30),
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.JWTSecret == "your-secret-key-change-in-production" && c.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}
	return nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSL,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func (c *Config) GetQdrantAddr() string {
	return fmt.Sprintf("http://%s:%s", c.QdrantHost, c.QdrantPort)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}
