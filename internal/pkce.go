package internal

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

var Plain Params

const (
	// code challenge methods defined as https://tools.ietf.org/html/rfc7636#section-4.3
	MethodS256 = "S256"
)

func NewState() (string, error) {
	b, err := Random32()
	if err != nil {
		return "", fmt.Errorf("could not generate a random: %w", err)
	}
	return Base64URLEncode(b), nil
}

func NewNonce() (string, error) {
	b, err := Random32()
	if err != nil {
		return "", fmt.Errorf("could not generate a random: %w", err)
	}
	return Base64URLEncode(b), nil
}

// Params represents a set of the PKCE parameters.
type Params struct {
	CodeChallenge       string
	CodeChallengeMethod string
	CodeVerifier        string
}

func (p Params) IsZero() bool {
	return p == Params{}
}

// New returns a parameters supported by the provider.
// You need to pass the code challenge methods defined in RFC7636.
// It returns Plain if no method is available.
func NewParam(methods []string) (Params, error) {
	for _, method := range methods {
		if method == MethodS256 {
			return NewS256()
		}
	}
	return Plain, nil
}

// NewS256 generates a parameters for S256.
func NewS256() (Params, error) {
	b, err := Random32()
	if err != nil {
		return Plain, fmt.Errorf("could not generate a random: %w", err)
	}
	return ComputeS256(b), nil
}

func Random32() ([]byte, error) {
	b := make([]byte, 32)
	if err := binary.Read(rand.Reader, binary.LittleEndian, b); err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}
	return b, nil
}

func ComputeS256(b []byte) Params {
	v := Base64URLEncode(b)
	s := sha256.New()
	_, _ = s.Write([]byte(v))
	return Params{
		CodeChallenge:       Base64URLEncode(s.Sum(nil)),
		CodeChallengeMethod: MethodS256,
		CodeVerifier:        v,
	}
}

func Base64URLEncode(b []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}
