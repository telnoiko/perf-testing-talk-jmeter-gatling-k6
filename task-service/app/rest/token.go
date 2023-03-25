package rest

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

// Token extracts a token from the Authorization header.
// The header is expected to match the format "Bearer XX", where "XX" is the
// JWT token.Ø
type Token struct{}

var (
	ErrNoTokenInRequest = errors.New("no token present in request")
)

// see github.com/golang-jwt/jwt/v4@v4.4.3/request/extractor.go:89

func (s *Token) ExtractToken(req *http.Request) (string, error) {
	tokenHeader := req.Header.Get("Authorization")
	// The usual convention is for "Bearer" to be title-cased. However, there's no
	// strict rule around this, and it's best to follow the robustness principle here.
	if tokenHeader == "" || !strings.HasPrefix(strings.ToLower(tokenHeader), "bearer ") {
		return "", ErrNoTokenInRequest
	}
	return tokenHeader[7:], nil
}

func (s *Token) GenerateJWT(email string) (string, error) {
	bits := make([]byte, 12)
	_, err := rand.Read(bits)
	if err != nil {
		panic(err)
	}

	claims := jwt.MapClaims{
		"email": email,
		"iss":   "task-service",
		"jti":   base64.StdEncoding.EncodeToString(bits),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(secretKey))
}
