package supervisor

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type Worker struct {
	ID     uuid.UUID
	Name   string
	Quit   chan struct{}
	Done   chan struct{}
	Jobs   <-chan Job
	Errors chan error
}

func NewWorker(name string, jobs <-chan Job) (*Worker, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Worker{
		ID:     id,
		Name:   name,
		Quit:   make(chan struct{}),
		Jobs:   jobs,
		Errors: make(chan error, 10),
	}, nil
}

func (w *Worker) Start(ctx context.Context) {
	slog.Info("Worker starting", "worker", w.Name)
	defer func() {
		close(w.Done)
		slog.Info("Worker stopped", "worker", w.Name)
	}()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Context cancelled, stopping worker", "worker", w.Name)
			return
		case <-w.Quit:
			slog.Info("Received quit signal, stopping worker", "worker", w.Name)
			return
		case job, ok := <-w.Jobs:
			if !ok {
				slog.Info("Jobs channel closed, stopping worker", "worker", w.Name)
				return
			}
			if err := job.Execute(ctx); err != nil {
				select {
				case w.Errors <- err:
				default:
					slog.Error("Worker error channel full", "worker", w.Name, "error", err)
				}
			}
		}
	}
}
