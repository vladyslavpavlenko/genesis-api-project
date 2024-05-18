package models

import (
	"gorm.io/gorm"
	"time"
)

const dbTimeout = time.Second * 3

var db *gorm.DB

// User is a GORM user model.
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"unique" json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Currency is a GORM currency model.
type Currency struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// Subscription is a GORM subscription model.
type Subscription struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	User             User      `gorm:"foreignKey:UserID" json:"-"`
	BaseCurrencyID   uint      `gorm:"not null;index" json:"base_currency_id"`
	BaseCurrency     Currency  `gorm:"foreignKey:BaseCurrencyID" json:"-"`
	TargetCurrencyID uint      `gorm:"not null;index" json:"target_currency_id"`
	TargetCurrency   Currency  `gorm:"foreignKey:TargetCurrencyID" json:"-"`
	CreatedAt        time.Time `json:"created_at"`
}

// Models is the type for this package. Any model that is included as a member in this type is available throughout
// the application, anywhere that the app variable is used, provided that the model is also added in the New function.
type Models struct {
	User         User
	Currency     Currency
	Subscription Subscription
}

// New is the function used to create an instance of the models package. It returns the type Model, which embeds
// all the types we want to be available to our application.
func New(dbPool *gorm.DB) Models {
	db = dbPool

	return Models{
		User:         User{},
		Currency:     Currency{},
		Subscription: Subscription{},
	}
}

func (c *Currency) GetIDbyCode(code string) (uint, error) {
	var currency Currency
	err := db.Where("code = ?", code).First(&currency).Error
	if err != nil {
		return 0, err
	}

	return currency.ID, nil
}
