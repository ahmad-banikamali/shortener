package outbound_test

import (
	"testing"

	"github.com/ahmad-banikamali/internal/adapter/outbound"
	"github.com/ahmad-banikamali/internal/adapter/outbound/storetest"
	"github.com/ahmad-banikamali/internal/core/port"
)

func TestMemoryStore(t *testing.T) {
	storetest.Run(t, func(t *testing.T) port.LinkStore {
		return outbound.NewMemoryStore()
	})
}
