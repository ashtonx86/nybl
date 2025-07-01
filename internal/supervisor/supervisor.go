package supervisor

import (
	"context"
	"log/slog"
	"sync"
)

// Oversees workers and dependencies
type Supervisor struct {
	Workers []Worker
	Singletons map[string]Singleton
}

func New() *Supervisor {
	return &Supervisor{
		Workers: make([]Worker, 0),
		Singletons: make(map[string]Singleton, 0),
	}
}

func (s *Supervisor) AddWorker(w Worker) {
	s.Workers = append(s.Workers, w)
}

func (s *Supervisor) AddSingleton(id string, singleton Singleton) {
	s.Singletons[id] = singleton
}

func (s *Supervisor) StartWorkers(ctx context.Context) {
	slog.Info("Starting workers...", "totalWorkers", len(s.Workers))
	var wg sync.WaitGroup 
	wg.Add(len(s.Workers))

	for _, w := range s.Workers {
		go func(worker Worker) {
			defer wg.Done()
			worker.Start(ctx)
		}(w)
	}

	wg.Wait()
}

func GetSingletonAs[T Singleton](s *Supervisor, id string) (T, bool) {
	singleton, ok := s.Singletons[id]
	if !ok {
		var zero T 
		return zero, false
	}
	casted, ok := singleton.(T)
	return casted, ok
}

func (s *Supervisor) InitSingletons(ctx context.Context) {
	slog.Info("Initializing and starting singletons...", "totalSingletons", len(s.Singletons))
	for _, singleton := range s.Singletons {
		if err := singleton.Init(ctx); err != nil {
			slog.Error("Error while starting singleton", "err", err)
		}
	}
}

func (s *Supervisor) StopSingletons(ctx context.Context) {
	for _, singleton := range s.Singletons {
		if err := singleton.Stop(ctx); err != nil {
			slog.Error("Error while stopping singleton", "err", err)
		}
	}
}

func (s *Supervisor) StopWorkers() {
	for _, w := range s.Workers {
		close(w.Quit)
	}
}
