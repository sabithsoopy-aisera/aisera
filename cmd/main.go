package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	aisera "go.aisera.cloud"
	"golang.org/x/exp/slog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGINT)
	defer cancel()

	err := callAsieraEndpoints(ctx)
	if err != nil {
		slog.Error("error calling aisera endpoints", err)
		os.Exit(1)
	}
}

func callAsieraEndpoints(ctx context.Context) error {
	aiserOffering, err := aisera.Login(ctx, aisera.LoginRequest{
		Username: os.Getenv("AISERA_USERNAME"),
		Password: os.Getenv("AISERA_PASSWORD"),
	})
	if err != nil {
		return fmt.Errorf("error logging in: %w", err)
	}
	bots, err := aiserOffering.Bots(ctx)
	if err != nil {
		return fmt.Errorf("error getting bots: %w", err)
	}
	log.Printf("bots: %+v", bots)
	return nil
}
