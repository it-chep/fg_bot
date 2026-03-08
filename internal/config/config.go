package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BotConfig
	DBConfig
	port         string
	errorAdminID int64
}

type BotConfig struct {
	token      string
	webhookURL string
	useWebhook bool
	isActive   bool
}

type DBConfig struct {
	dbUser     string
	dbPassword string
	dbHost     string
	dbName     string
}

func (c *Config) Token() string {
	return c.token
}

func (c *Config) WebhookURL() string {
	return c.webhookURL
}

func (c *Config) UseWebhook() bool {
	return c.useWebhook
}

func (c *Config) BotIsActive() bool {
	return c.isActive
}

func (c *Config) Port() string {
	return c.port
}

func (c *Config) ErrorAdminID() int64 {
	return c.errorAdminID
}

func (c *Config) DBUser() string {
	return c.dbUser
}

func (c *Config) DBPassword() string {
	return c.dbPassword
}

func (c *Config) DBHost() string {
	return c.dbHost
}

func (c *Config) DBName() string {
	return c.dbName
}

func (c *Config) DatabaseURL() string {
	return "postgres://" + c.dbUser + ":" + c.dbPassword + "@" + c.dbHost + ":5432/" + c.dbName
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		BotConfig: BotConfig{
			webhookURL: os.Getenv("WEBHOOK_URL"),
			token:      os.Getenv("BOT_TOKEN"),
			useWebhook: os.Getenv("USE_WEBHOOK") == "true",
			isActive:   os.Getenv("BOT_IS_ACTIVE") == "true",
		},
		DBConfig: DBConfig{
			dbUser:     os.Getenv("DB_USER"),
			dbPassword: os.Getenv("DB_PASSWORD"),
			dbHost:     os.Getenv("DB_HOST"),
			dbName:     os.Getenv("DB_NAME"),
		},
		port:         port,
		errorAdminID: parseInt64(os.Getenv("ERROR_ADMIN_ID")),
	}
}

func parseInt64(raw string) int64 {
	if raw == "" {
		return 0
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		log.Printf("invalid ERROR_ADMIN_ID value: %q", raw)
		return 0
	}
	return v
}
