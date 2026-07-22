package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
	tokenIssuer     = "code-colosseum"
)

type TokenManager struct {
	secret []byte
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{
		secret: []byte(secret),
	}
}

func (m *TokenManager) GenerateAccessToken(userID string) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    tokenIssuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signedToken, nil
}

func GenerateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (m *TokenManager) ValidateAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(tokenIssuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return "", fmt.Errorf("validate access token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid access token")
	}

	if claims.Subject == "" {
		return "", errors.New("access token has no user ID")
	}

	return claims.Subject, nil
}
