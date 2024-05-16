package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"net/http"
	"time"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// subscriptionBody is the email subscription request body structure.
type subscriptionBody struct {
	Email              string `json:"email"`
	BaseCurrencyCode   string `json:"base_currency_code"`
	TargetCurrencyCode string `json:"target_currency_code"`
}

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewTestRepo creates a new test repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// GetRate gets the current USD to UAH exchange rate.
func (m *Repository) GetRate(w http.ResponseWriter, r *http.Request) {

}

// Subscribe handles email subscription.
func (m *Repository) Subscribe(w http.ResponseWriter, r *http.Request) {
	var body subscriptionBody

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if err != nil {
		_ = m.errorJSON(w, errors.New("failed to read body"))
		return
	}

	if !validateEmail(body.Email) {
		_ = m.errorJSON(w, errors.New("email is invalid"))
		return
	}

	// Create a new user model
	user := models.User{
		Email:     body.Email,
		CreatedAt: time.Now(),
	}

	// Add user to the database
	result := m.App.DB.Create(&user)
	if result.Error != nil {
		_ = m.errorJSON(w, fmt.Errorf("error creating user"), http.StatusInternalServerError)
		return
	}

	// Get IDs of the currencies by their codes
	var baseCurrency = models.Currency{
		Code: body.BaseCurrencyCode,
	}

	baseCurrencyID, err := m.App.Models.Currency.GetIDbyCode(baseCurrency.Code)
	if err != nil {
		_ = m.errorJSON(w, fmt.Errorf("error retrieving base currency"))
		return
	}

	var targetCurrency = models.Currency{
		Code: body.TargetCurrencyCode,
	}

	targetCurrencyID, err := m.App.Models.Currency.GetIDbyCode(targetCurrency.Code)
	if err != nil {
		_ = m.errorJSON(w, fmt.Errorf("error retrieving target currency"))
		return
	}

	// Create a new subscription model
	subscription := models.Subscription{
		UserID:           user.ID,
		BaseCurrencyID:   baseCurrencyID,
		TargetCurrencyID: targetCurrencyID,
	}

	// Add subscription to the database
	result = m.App.DB.Create(&subscription)
	if result.Error != nil {
		_ = m.errorJSON(w, fmt.Errorf("error creating subscription"), http.StatusInternalServerError)
		return
	}

	// Send response
	payload := jsonResponse{
		Error:   false,
		Message: "subscribed",
	}

	_ = m.writeJSON(w, http.StatusAccepted, payload)
}

// SendEmails handles sending emails to all the subscribed emails.
func (m *Repository) SendEmails(w http.ResponseWriter, r *http.Request) {

}
