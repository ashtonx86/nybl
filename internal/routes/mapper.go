package routes

import (
	v1Router "github.com/ashtonx86/nybl/internal/routes/api/v1"
	"github.com/ashtonx86/nybl/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

func MapAllRoutes(su *supervisor.Supervisor, app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "1")
		return c.Next()
	})

	v1Router.MapAccountsRoute(su, v1.(*fiber.Group))
}
