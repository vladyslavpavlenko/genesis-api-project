package mailer

import (
	"github.com/robfig/cron/v3"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"gorm.io/gorm"
)

// ScheduleEmails schedules the email sending task to run at 10 AM every day.
func ScheduleEmails(emailConfig config.EmailConfig, db *gorm.DB) {
	c := cron.New()

	c.AddFunc("0 10 * * *", func() {
		SendEmails(emailConfig, db)
	})

	c.Start()
}
