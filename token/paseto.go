package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

var (
	ErrInvalidKeySize     = fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	errInvalidPasetoToken = fmt.Errorf("error decrypting paseto token: %s", ErrInvalidToken)
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new Paseto maker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, ErrInvalidKeySize
	}

	// maker := &PasetoMaker{
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil

	// return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (mkr *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewTokenPayload(username, duration)
	if err != nil {
		return "", fmt.Errorf("error creating paseto payload: %s", err)
	}

	return mkr.paseto.Encrypt(mkr.symmetricKey, payload, nil)
}

// VerifyToken checks if the token is valid or not
func (mkr *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := mkr.paseto.Decrypt(token, mkr.symmetricKey, payload, nil)
	if err != nil {
		return nil, errInvalidPasetoToken
	}

	err = payload.Valid()
	if err != nil {
		// return nil, fmt.Errorf("error validating paseto token: %s", err)
		return nil, err
	}

	return payload, nil
}
