package main

import (
	"context"

	"github.com/ashtonx86/nybl/internal/server"
)

func main() {
	server := server.NewWebServer()
	server.StartChildServices()
	server.MapRoutes()
	defer server.GracefulShutdown(context.Background())

	server.App.Listen(":3000")
}
