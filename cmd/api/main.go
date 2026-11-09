package main

import (
	"log"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
