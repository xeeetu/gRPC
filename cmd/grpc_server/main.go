package main

import (
	"context"
	"flag"
	"log"

	"github.com/xeeetu/gRPC/internal/app"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
