package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

type AuthConfig struct {
	codeVerifier string;
	codeChallenge string;
	state string;
}


func generateRandomString(length int) string {
	b := make([]byte, length)
	_,err := rand.Read(b)
	if err != nil {
		panic("rand.Read failed")
	}
  return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(codeVerifier string) string {
	h := sha256.New()
  h.Write([]byte(codeVerifier))
  return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func NewAuthConfig() *AuthConfig {
	codeVerifier := generateRandomString(32)
	state := generateRandomString(32)
  codeChallenge := generateCodeChallenge(codeVerifier)

  return &AuthConfig{
    codeVerifier: codeVerifier,
    codeChallenge: codeChallenge,
    state: state,
  }
}
