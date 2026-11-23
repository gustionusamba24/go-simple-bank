package token

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gustionusamba24/simple-bank/util"
	"github.com/stretchr/testify/require"
)

// failingReader is an io.Reader that always fails, used to force uuid.NewRandom to error.
type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errors.New("random source failure")
}

// withFailingUUID makes uuid generation fail for the duration of fn.
func withFailingUUID(t *testing.T, fn func()) {
	uuid.SetRand(failingReader{})
	t.Cleanup(func() { uuid.SetRand(nil) })
	fn()
}

func TestCreateTokenPayloadError(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	pasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	makers := map[string]Maker{"jwt": jwtMaker, "paseto": pasetoMaker}

	for name, maker := range makers {
		t.Run(name, func(t *testing.T) {
			withFailingUUID(t, func() {
				token, payload, err := maker.CreateToken(util.RandomOwner(), time.Minute)
				require.Error(t, err)
				require.Empty(t, token)
				require.Nil(t, payload)
			})
		})
	}
}
