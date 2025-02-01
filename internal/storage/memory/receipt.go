package memory

import (
	"fmt"
	"sync"
)

// ReceiptStorage defines the interface for storing receipt points
type ReceiptStorage interface {
	SavePoints(id string, points int) error
	GetPoints(id string) (int, error)
}

// MemoryStorage implements ReceiptStorage using in-memory map
type MemoryStorage struct {
	points map[string]int // map of receipt id to points
	mu     sync.RWMutex   // protects concurrent access to points
}

// NewMemoryStorage creates a new MemoryStorage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		points: make(map[string]int),
	}
}

// SavePoints saves the points for a receipt id
func (s *MemoryStorage) SavePoints(id string, points int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.points[id] = points
	return nil
}

// GetPoints returns the points for a receipt id
func (s *MemoryStorage) GetPoints(id string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	points, exists := s.points[id]

	if !exists {
		return 0, fmt.Errorf("points not found for id: %s", id)
	}

	return points, nil
}
