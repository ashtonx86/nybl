package supervisor

import "context"

const (
	StatusStarting = "starting"
	StatusRunning  = "running"
	StatusStopping = "stopping"
	StatusStopped  = "stopped"
	StatusError    = "error"
)

// Long lived, single instance service.
type Singleton interface {
	Init(ctx context.Context) error
	Stop(ctx context.Context) error
}