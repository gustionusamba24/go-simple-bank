package token

import (
	"testing"
	"time"

	"github.com/gustionusamba24/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestNewPasetoMakerInvalidKeySize(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(31))
	require.Error(t, err)
	require.Nil(t, maker)

	maker, err = NewPasetoMaker(util.RandomString(33))
	require.Error(t, err)
	require.Nil(t, maker)
}

func TestInvalidPasetoTokenWrongKey(t *testing.T) {
	maker1, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	maker2, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, _, err := maker1.CreateToken(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	payload, err := maker2.VerifyToken(token)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoTokenMalformed(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken("not-a-valid-token")
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
