package outbound_test

import (
	"path/filepath"
	"testing"

	"github.com/ahmad-banikamali/shortener/internal/adapter/outbound"
	"github.com/ahmad-banikamali/shortener/internal/adapter/outbound/storetest"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

func TestSQLiteStore(t *testing.T) {
	storetest.Run(t, func(t *testing.T) port.LinkStore {
		s, err := outbound.NewSQLiteStore(filepath.Join(t.TempDir(), "test.db"))
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { s.Close() })
		return s
	})
}
