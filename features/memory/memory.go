package memory

import (
    "net/http"
    "github.com/brnsampson/echopilot/features/memory/internal/rest_handler"
    "github.com/brnsampson/echopilot/features/memory/records"
    "github.com/charmbracelet/log"
)

// Interfaces for strategy pattern
type Store[T any] interface {
    Create(T) (T, error)
    Update(T) (T, error)
    Patch(T) (T, error)
    Get(T) T
    List(T, int, int) []T
    Delete(T) (T, error)
}

// Service
func NewFeature(logger *log.Logger) *Feature {
    memories := records.NewMemoryStore()

	return &Feature{ logger.With("package", "memory"), memories}
}

type Feature struct{
    logger *log.Logger
    store *records.MemoryStore
}

func (f *Feature) GetHandler() (string, http.Handler) {
    return "/memory", rest_handler.NewRestHandler(f.store)
}
