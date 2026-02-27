package lsm_test

import (
	"testing"

	"github.com/aaronshifman/lsm-cause-i-can/pkg/lsm"
	"github.com/stretchr/testify/require"
)

func TestMemoryCache_Get(t *testing.T) {
	tests := []struct {
		name      string
		setup     map[string]int
		key       string
		wantVal   int
		wantFound bool
	}{
		{
			name:      "key exists",
			setup:     map[string]int{"a": 1},
			key:       "a",
			wantVal:   1,
			wantFound: true,
		},
		{
			name:      "key missing",
			setup:     map[string]int{"a": 1},
			key:       "b",
			wantVal:   0,
			wantFound: false,
		},
		{
			name:      "empty cache",
			setup:     map[string]int{},
			key:       "a",
			wantVal:   0,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := lsm.NewMemoryCache[string, int]()
			for k, v := range tt.setup {
				c.Put(k, v)
			}

			got, ok := c.Get(tt.key)
			require.Equal(t, tt.wantFound, ok)
			require.Equal(t, tt.wantVal, got)
		})
	}
}

func TestMemoryCache_Put(t *testing.T) {
	tests := []struct {
		name    string
		puts    []struct{ key, val string }
		getKey  string
		wantVal string
	}{
		{
			name:    "insert new key",
			puts:    []struct{ key, val string }{{"x", "hello"}},
			getKey:  "x",
			wantVal: "hello",
		},
		{
			name: "overwrite existing key",
			puts: []struct{ key, val string }{
				{"x", "first"},
				{"x", "second"},
			},
			getKey:  "x",
			wantVal: "second",
		},
		{
			name: "multiple keys are independent",
			puts: []struct{ key, val string }{
				{"a", "alpha"},
				{"b", "beta"},
			},
			getKey:  "a",
			wantVal: "alpha",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := lsm.NewMemoryCache[string, string]()
			for _, p := range tt.puts {
				c.Put(p.key, p.val)
			}

			got, ok := c.Get(tt.getKey)
			require.True(t, ok)
			require.Equal(t, tt.wantVal, got)
		})
	}
}
