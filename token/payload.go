package token

import (
	"errors"
	"time"

	// 这里需要 go get github.com/gofrs/uuid
	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken funciton 
var (
	ErrorExpireToken = errors.New("token has expired")
	ErrorInvalidToken = errors.New("token is invalid")
)
// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks if the token payload is valid or not 
func (payload *Payload) Valid() error {
	// 检查 payload 是否过期
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpireToken
	}
	return nil
}