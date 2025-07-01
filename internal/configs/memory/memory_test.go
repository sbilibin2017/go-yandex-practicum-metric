package memory

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemory_Default(t *testing.T) {
	mem := NewMemory[string, int]()

	assert.NotNil(t, mem.Mu, "Expected default mutex to be initialized")
	assert.NotNil(t, mem.Data, "Expected default data map to be initialized")
	assert.Empty(t, mem.Data, "Expected default data map to be empty")
}

func TestNewMemory_WithMutex(t *testing.T) {
	mu := &sync.RWMutex{}
	mem := NewMemory[string, int](
		WithMutex[string, int](mu),
	)

	assert.Equal(t, mu, mem.Mu)
	assert.NotNil(t, mem.Data, "Expected data to be initialized by default")
}

func TestNewMemory_WithData(t *testing.T) {
	initialData := map[string]int{"a": 1, "b": 2}
	mem := NewMemory(
		WithData[string, int](initialData),
	)

	assert.Equal(t, initialData, mem.Data)
	assert.NotNil(t, mem.Mu, "Expected mutex to be initialized by default")
}

func TestNewMemory_WithAllOptions(t *testing.T) {
	mu := &sync.RWMutex{}
	data := map[string]int{"x": 42}

	mem := NewMemory(
		WithMutex[string, int](mu),
		WithData[string, int](data),
	)

	assert.Equal(t, mu, mem.Mu)
	assert.Equal(t, data, mem.Data)
}

func TestMemory_ConcurrentAccess(t *testing.T) {
	mu := &sync.RWMutex{}
	data := map[string]int{}
	mem := NewMemory(
		WithMutex[string, int](mu),
		WithData[string, int](data),
	)

	require.NotNil(t, mem.Mu)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		mem.Mu.Lock()
		defer mem.Mu.Unlock()
		mem.Data["key"] = 100
	}()

	go func() {
		defer wg.Done()
		mem.Mu.RLock()
		defer mem.Mu.RUnlock()
		_ = mem.Data["key"]
	}()

	wg.Wait()
	assert.Equal(t, 100, mem.Data["key"])
}
