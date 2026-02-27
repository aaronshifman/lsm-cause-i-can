package lsm

// MemoryCache is the in-memory component of the LSM tree
type MemoryCache[K comparable, V any] struct {
	// mem is the in-memory K/V store of the LSM tree
	mem map[K]V
}

// NewMemoryCache creates the empty K/V store
func NewMemoryCache[K comparable, V any]() *MemoryCache[K, V] {
	m := make(map[K]V)
	return &MemoryCache[K, V]{
		mem: m,
	}
}

// Get looks ok a key of type K in the cache and returns either the key or a zero and
// the success of the lookup
func (m *MemoryCache[K, V]) Get(key K) (V, bool) {
	k, ok := m.mem[key]
	var zero V
	if !ok {
		return zero, false
	}

	return k, true
}

// Put upserts the value into the key
// Squashes existing key value if exists
func (m *MemoryCache[K, V]) Put(key K, val V) {
	m.mem[key] = val
}
