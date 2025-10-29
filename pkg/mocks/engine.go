// ---------------------------
// pkg/mocks/engine.go
// ---------------------------
package mocks

import "fmt"

// Simple in-memory mock store. In production you'd design a stricter schema.
type MockStore struct {
	values map[string]map[string]string // resource -> attributes
}

func NewMockStore() *MockStore {
	return &MockStore{values: make(map[string]map[string]string)}
}

func (s *MockStore) Add(resource string, attrs map[string]string) {
	s.values[resource] = attrs
}

func (s *MockStore) Get(resource string) (map[string]string, bool) {
	v, ok := s.values[resource]
	return v, ok
}

func (s *MockStore) String() string {
	return fmt.Sprintf("%v", s.values)
}
