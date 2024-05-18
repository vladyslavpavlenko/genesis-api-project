package mailer

import (
	"fmt"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/rate"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"log"
	"sync"
)

// EmailMessage holds the email message data.
type EmailMessage struct {
	From    string
	To      string
	Subject string
	Body    string
}

// RateUpdate holds the exchange rate update data.
type RateUpdate struct {
	Base   string
	Target string
	Price  string
}

// SendEmail sends an email using the provided configuration and message data.
func SendEmail(wg *sync.WaitGroup, emailConfig config.EmailConfig, recipient string, update RateUpdate) {
	defer wg.Done()

	msg := EmailMessage{
		From:    emailConfig.Email,
		To:      recipient,
		Subject: fmt.Sprintf("%s to %s Exchange Rate", update.Base, update.Target),
		Body:    fmt.Sprintf("The current exchange rate for %s to %s is %s.", update.Base, update.Target, update.Price),
	}

	m := gomail.NewMessage()
	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/plain", msg.Body)

	d := gomail.NewDialer("smtp.gmail.com", 587, emailConfig.Email, emailConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Could not send email to %s: %v", msg.To, err)
	} else {
		log.Printf("Email sent successfully to %s!", msg.To)
	}
}

// SendEmails sends emails to all subscribers in the database.
func SendEmails(emailConfig config.EmailConfig, db *gorm.DB) {
	var subscriptions []models.Subscription
	db.Preload("User").Preload("BaseCurrency").Preload("TargetCurrency").Find(&subscriptions)

	var wg sync.WaitGroup
	for _, subscription := range subscriptions {
		wg.Add(1)

		base := subscription.BaseCurrency.Code
		target := subscription.TargetCurrency.Code
		price, err := rate.GetRate(base, target)
		if err != nil {
			log.Printf("Failed to retrieve rate for %s to %s: %v", base, target, err)
			wg.Done()
			continue
		}

		var update = RateUpdate{
			Base:   base,
			Target: target,
			Price:  price,
		}

		go SendEmail(&wg, emailConfig, subscription.User.Email, update)
	}
	wg.Wait()
}
