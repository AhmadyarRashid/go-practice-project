package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
	RateLimit RateLimitConfig
	CORS     CORSConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name  string
	Env   string
	Port  string
	Debug bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret             string
	ExpiryHours        int
	RefreshExpiryHours int
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int
	Duration time.Duration
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Read .env file if it exists
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if .env doesn't exist, we can use environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Set defaults
	setDefaults()

	config := &Config{
		App: AppConfig{
			Name:  viper.GetString("APP_NAME"),
			Env:   viper.GetString("APP_ENV"),
			Port:  viper.GetString("APP_PORT"),
			Debug: viper.GetBool("APP_DEBUG"),
		},
		Database: DatabaseConfig{
			Driver:   viper.GetString("DB_DRIVER"),
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			Name:     viper.GetString("DB_NAME"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
		},
		JWT: JWTConfig{
			Secret:             viper.GetString("JWT_SECRET"),
			ExpiryHours:        viper.GetInt("JWT_EXPIRY_HOURS"),
			RefreshExpiryHours: viper.GetInt("JWT_REFRESH_EXPIRY_HOURS"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
		RateLimit: RateLimitConfig{
			Requests: viper.GetInt("RATE_LIMIT_REQUESTS"),
			Duration: viper.GetDuration("RATE_LIMIT_DURATION"),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ","),
			AllowedMethods: strings.Split(viper.GetString("CORS_ALLOWED_METHODS"), ","),
			AllowedHeaders: strings.Split(viper.GetString("CORS_ALLOWED_HEADERS"), ","),
		},
	}

	// Validate required configurations
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// setDefaults sets default values for configuration
func setDefaults() {
	viper.SetDefault("APP_NAME", "go-enterprise-api")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_DEBUG", true)

	viper.SetDefault("DB_DRIVER", "sqlite")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_NAME", "enterprise.db")
	viper.SetDefault("DB_SSL_MODE", "disable")

	viper.SetDefault("JWT_EXPIRY_HOURS", 24)
	viper.SetDefault("JWT_REFRESH_EXPIRY_HOURS", 168)

	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "json")

	viper.SetDefault("RATE_LIMIT_REQUESTS", 100)
	viper.SetDefault("RATE_LIMIT_DURATION", "1m")

	viper.SetDefault("CORS_ALLOWED_ORIGINS", "*")
	viper.SetDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS")
	viper.SetDefault("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Authorization")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if c.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	return nil
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	switch c.Database.Driver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host,
			c.Database.Port,
			c.Database.User,
			c.Database.Password,
			c.Database.Name,
			c.Database.SSLMode,
		)
	case "sqlite":
		return c.Database.Name
	default:
		return ""
	}
}
