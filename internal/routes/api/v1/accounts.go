package v1

import (
	"github.com/ashtonx86/nybl/internal/handlers"
	"github.com/ashtonx86/nybl/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

func MapAccountsRoute(su *supervisor.Supervisor, r *fiber.Group) {
	g := r.Group("/accounts")
	handler := handlers.NewAccountsHandler(su)

	g.Get("/", handler.Get)
	g.Post("/", handler.Register)
}
