package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/mailer"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/rate"
	"net/http"
	"strings"
	"time"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type rateResponse struct {
	BaseCurrencyCode   string `json:"base_currency_code"`
	TargetCurrencyCode string `json:"target_currency_code"`
	Price              string `json:"price"`
}

// subscriptionBody is the email subscription request body structure.
type subscriptionBody struct {
	Email string `json:"email"`
	// BaseCurrencyCode   string `json:"base_currency_code"`
	// TargetCurrencyCode string `json:"target_currency_code"`
}

// GetRate gets the current USD to UAH exchange rate.
func (m *Repository) GetRate(w http.ResponseWriter, r *http.Request) {
	price, err := rate.GetRate("USD", "UAH")
	if err != nil {
		_ = m.errorJSON(w, errors.New("error calling Coinbase API"), http.StatusBadRequest) // http.StatusServiceUnavailable
		return
	}

	rateResp := rateResponse{
		BaseCurrencyCode:   "USD",
		TargetCurrencyCode: "UAH",
		Price:              price,
	}

	// Send response
	payload := jsonResponse{
		Error: false,
		Data:  rateResp,
	}

	_ = m.writeJSON(w, http.StatusOK, payload)
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
		if strings.Contains(result.Error.Error(), "duplicate") ||
			strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			_ = m.errorJSON(w, fmt.Errorf("already subscribed"), http.StatusConflict)
			return
		}

		_ = m.errorJSON(w, fmt.Errorf("error creating user"), http.StatusInternalServerError)
		return
	}

	// Get IDs of the currencies by their codes
	var baseCurrency = models.Currency{
		Code: "USD", // body.BaseCurrencyCode,
	}

	baseCurrencyID, err := m.App.Models.Currency.GetIDbyCode(baseCurrency.Code)
	if err != nil {
		_ = m.errorJSON(w, fmt.Errorf("error retrieving base currency"))
		return
	}

	var targetCurrency = models.Currency{
		Code: "UAH", // body.BaseCurrencyCode,
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
		if strings.Contains(result.Error.Error(), "duplicate") ||
			strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			_ = m.errorJSON(w, fmt.Errorf("already subscribed"), http.StatusConflict)
			return
		}

		_ = m.errorJSON(w, fmt.Errorf("error creating subscription"), http.StatusInternalServerError)
		return
	}

	// Send response
	payload := jsonResponse{
		Error:   false,
		Message: "subscribed",
	}

	_ = m.writeJSON(w, http.StatusOK, payload)
}

// SendEmails handles sending emails to all the subscribed emails.
func (m *Repository) SendEmails(w http.ResponseWriter, r *http.Request) {
	mailer.SendEmails(m.App.EmailConfig, m.App.DB)
}
