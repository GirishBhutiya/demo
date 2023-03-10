package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	synmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(synmetricKey string) (Maker, error) {
	if len(synmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		synmetricKey: []byte(synmetricKey),
	}
	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.synmetricKey, payload, nil)
}

// VerifyToken check if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.synmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()

	if err != nil {
		return nil, err
	}
	return payload, nil
}
