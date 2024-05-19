package handlers

import (
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/mail"
	"time"
)

func (m *Repository) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1 MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func (m *Repository) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON writes the error JSON. If status is not provided, the default http.StatusBadRequest in sent.
func (m *Repository) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return m.writeJSON(w, statusCode, payload)
}

// createUser creates a new user and saves it to the database.
func (m *Repository) createUser(email string) (*models.User, error) {
	user := models.User{
		Email:     email,
		CreatedAt: time.Now(),
	}
	result := m.App.DB.Create(&user)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return nil, gorm.ErrDuplicatedKey
		}
		return nil, result.Error
	}
	return &user, nil
}

// createSubscription creates a new subscription and saves it to the database.
func (m *Repository) createSubscription(userID, baseCurrencyID, targetCurrencyID uint) error {
	subscription := models.Subscription{
		UserID:           userID,
		BaseCurrencyID:   baseCurrencyID,
		TargetCurrencyID: targetCurrencyID,
	}
	result := m.App.DB.Create(&subscription)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return gorm.ErrDuplicatedKey
		}
		return result.Error
	}
	return nil
}

// getCurrencyID retrieves the ID of a currency by its code.
func (m *Repository) getCurrencyID(code string) (uint, error) {
	currencyID, err := m.App.Models.Currency.GetIDbyCode(code)
	if err != nil {
		return 0, err
	}
	return currencyID, nil
}

// validateEmail validates the email address.
func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
