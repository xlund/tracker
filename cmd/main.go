package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/xlund/tracker/internal/app"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	ret := app.Start(ctx)
	os.Exit(ret)
}
