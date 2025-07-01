package connectionschemas

import (
	"github.com/ashtonx86/nybl/internal/schemas"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Conn   *websocket.Conn
	Author schemas.Account
}
