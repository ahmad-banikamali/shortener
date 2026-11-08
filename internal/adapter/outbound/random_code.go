package outbound

import (
	"crypto/rand"
	"math/big"

	"github.com/ahmad-banikamali/internal/core/port"
)

var _ port.CodeGenerator = RandomCodeGenerator{}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RandomCodeGenerator struct {
	Length int
}

func (g RandomCodeGenerator) Generate() (string, error) {
	_max := big.NewInt(int64(len(alphabet)))
	b := make([]byte, g.Length)
	for i := range b {
		n, err := rand.Int(rand.Reader, _max)
		if err != nil {
			return "", err
		}
		b[i] = alphabet[n.Int64()]
	}
	return string(b), nil
}
