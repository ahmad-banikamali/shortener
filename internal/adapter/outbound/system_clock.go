package outbound

import (
	"time"

	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

var _ port.Clock = SystemClock{}

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now() }
