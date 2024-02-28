package main

import (
	"speedstar/internal/services"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	go services.RunScheduler()

	// Block main thread from completing
	for {
		time.Sleep(1000)
	}
}
