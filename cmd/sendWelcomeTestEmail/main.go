package main

import (
	"log"

	"github.com/bmkersey/Go-SaaSy/internal/email"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sender := email.NewEmailSender()

	data := struct {
		Name string
	}{
		Name: "Test User",
	}

	err = sender.SendTemplate("bmkersey@gmail.com", "Welcome to SaaSy!", "welcome.html.tmpl", data)
	if err != nil {
		log.Fatalf("Error sending test email: %v", err)
	}

	log.Println("Test email sent successfully.")
}
