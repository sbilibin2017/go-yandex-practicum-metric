package memory

import "sync"

// Memory is a generic struct for storing metric data with concurrency support
type Memory[K comparable, V any] struct {
	Mu   *sync.RWMutex
	Data map[K]V
}

// Opt is a functional option type for configuring Memory
type Opt[K comparable, V any] func(*Memory[K, V])

// WithMutex allows setting a custom RWMutex (as a pointer)
func WithMutex[K comparable, V any](mu *sync.RWMutex) Opt[K, V] {
	return func(m *Memory[K, V]) {
		m.Mu = mu
	}
}

// WithData allows setting initial map data
func WithData[K comparable, V any](data map[K]V) Opt[K, V] {
	return func(m *Memory[K, V]) {
		m.Data = data
	}
}

// NewMemory constructs a Memory instance with optional configuration
func NewMemory[K comparable, V any](opts ...Opt[K, V]) *Memory[K, V] {
	m := &Memory[K, V]{
		Mu:   &sync.RWMutex{},
		Data: make(map[K]V),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}
