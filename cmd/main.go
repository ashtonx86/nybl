package main

import (
	"context"
	"flag"

	"github.com/ashtonx86/nybl/internal/server"
	"github.com/joho/godotenv"
)

var isDev = flag.Bool("dev", false, "-dev")

func main() {
	flag.Parse()

	if *isDev {
		if err := godotenv.Load(); err != nil {
			panic("Failed to load .env")
		}
	}

	server := server.NewWebServer()
	server.StartChildServices()
	server.MapRoutes()
	defer server.GracefulShutdown(context.Background())

	server.App.Listen(":3000")
}
