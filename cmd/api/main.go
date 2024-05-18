package main

import (
	"fmt"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/mailer"
	"log"
	"net/http"
)

const webPort = 8080

var app config.AppConfig

func main() {
	err := setup(&app)
	if err != nil {
		log.Fatal()
	}

	mailer.ScheduleEmails(app.EmailConfig, app.DB)

	log.Printf("Running on port %d", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
