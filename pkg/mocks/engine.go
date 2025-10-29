// ---------------------------
// pkg/mocks/engine.go
// ---------------------------
package mocks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

// InjectMocks writes .tf files into moduleDir to represent mocked resources.
// It generates a resource block per mock with attributes as constants.
func InjectMocks(moduleDir string, store *MockStore) error {
	for res, attrs := range store.values {
		// res is like "aws_vpc.main" -> type=aws_vpc, name=main
		parts := strings.Split(res, ".")
		if len(parts) != 2 {
			return fmt.Errorf("invalid mock resource format: %s", res)
		}
		typeName := parts[0]
		resName := parts[1]
		filename := filepath.Join(moduleDir, fmt.Sprintf("mock_%s_%s.tf", typeName, resName))
		f, err := os.Create(filename)
		if err != nil { return err }
		defer f.Close()
		fmt.Fprintf(f, "resource \"%s\" \"%s\" {", typeName, resName)
		for k, v := range attrs {
			// write attribute as string literal
			fmt.Fprintf(f, "  %s = \"%s\"", k, v)
		}
		// If id is set, set it as well via computed attribute 'id' isn't directly settable,
		// but many resources accept tags/other attributes. We'll also set a local output for test reads.
		fmt.Fprintf(f, "} output \"mock_%s_%s_attrs\" {value = {", typeName, resName)
		for k := range attrs {
			fmt.Fprintf(f, "    %s = resource.%s.%s.%s", k, typeName, resName, k)
		}
		fmt.Fprintf(f, "  }}")
	}
	return nil
}
