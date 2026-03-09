package token

import (
	"testing"
	"time"

	"github.com/Long4Changes/MySimpleBank/db/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	// 拿到 payload 之后一项一项检查
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}


func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	// 将过期时间设置为当前时间，人为构造立即过期的 token
	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpireToken.Error())
	require.Nil(t, payload)
}

// 检查使用 "alg" = None 的无效 token
// Only allow 'none' alg type if UnsafeAllowNoneSignatureType is specified as the key
// Key must be UnsafeAllowNoneSignatureType to prevent accidentally
// accepting 'none' signing method
func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorInvalidToken.Error())
	require.Nil(t, payload)
}