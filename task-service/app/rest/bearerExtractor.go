package rest

import (
	"errors"
	"net/http"
	"strings"
)

// BearerExtractor extracts a token from the Authorization header.
// The header is expected to match the format "Bearer XX", where "XX" is the
// JWT token.Ø
type BearerExtractor struct{}

var (
	ErrNoTokenInRequest = errors.New("no token present in request")
)

// see github.com/golang-jwt/jwt/v4@v4.4.3/request/extractor.go:89

func (e BearerExtractor) ExtractToken(req *http.Request) (string, error) {
	tokenHeader := req.Header.Get("Authorization")
	// The usual convention is for "Bearer" to be title-cased. However, there's no
	// strict rule around this, and it's best to follow the robustness principle here.
	if tokenHeader == "" || !strings.HasPrefix(strings.ToLower(tokenHeader), "bearer ") {
		return "", ErrNoTokenInRequest
	}
	return tokenHeader[7:], nil
}
