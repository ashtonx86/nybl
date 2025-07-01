package ws

import (
	"context"
	"sync"

	schemas "github.com/ashtonx86/nybl/internal/schemas/connection_schemas"
	"github.com/gofiber/fiber/v2"
)

type WebSocketGateway struct {
	mu                sync.RWMutex
	ActiveConnections map[string]schemas.Connection
}

func NewGateway() *WebSocketGateway {
	return &WebSocketGateway{
		ActiveConnections: make(map[string]schemas.Connection),
	}
}

func (g *WebSocketGateway) Promote(ctx context.Context, r *fiber.Request) {
	g.mu.RLock()
	defer g.mu.Unlock()
}
