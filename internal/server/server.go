package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/Auden747/nybl/internal/dependencies"
	"github.com/Auden747/nybl/internal/routes"
	"github.com/Auden747/nybl/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

type WebServer struct {
	App *fiber.App 
	Supervisor *supervisor.Supervisor
}

func NewWebServer() *WebServer {
	app := fiber.New()
	su := supervisor.New()

	su.AddSingleton("sqlite", dependencies.NewSQLiteSingleton())
	su.AddSingleton("email", dependencies.NewMailSingleton())

	return &WebServer{
		App: app,
		Supervisor: su,
	}
}

func (s *WebServer) StartChildServices() {
	slog.Info("Supervisor is whipping workers and singletons into action...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 220)
	defer cancel()

	s.Supervisor.InitSingletons(ctx)
	s.Supervisor.StartWorkers(ctx)
}

func (s *WebServer) MapRoutes() {
	routes.MapAllRoutes(s.Supervisor, s.App)
}

func (s *WebServer) GracefulShutdown(ctx context.Context) {
	err := s.App.ShutdownWithContext(ctx)
	if err != nil {
		slog.Error("Failed to gracefully shutdown server", "error", err)
	}

	s.Supervisor.StopSingletons(ctx)
	s.Supervisor.StopWorkers()
}