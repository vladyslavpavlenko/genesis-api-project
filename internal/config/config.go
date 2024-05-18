package config

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"gorm.io/gorm"
)

// AppConfig holds the application config.
type AppConfig struct {
	DB          *gorm.DB
	Models      models.Models
	EmailConfig EmailConfig
}

// EmailConfig holds the email configuration.
type EmailConfig struct {
	Email    string
	Password string
}
