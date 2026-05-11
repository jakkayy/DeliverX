package middleware

import (
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	Role string `json:"role"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

type JWTUtil struct {
	secret              []byte
	expirationMs        int64
	refreshExpirationMs int64
}

func NewJWTUtil(secret string, expirationMs, refreshExpirationMs int64) *JWTUtil {
	key, _ := base64.StdEncoding.DecodeString(secret)
	return &JWTUtil{secret: key, expirationMs: expirationMs, refreshExpirationMs: refreshExpirationMs}
}

func (j *JWTUtil) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := JWTClaims{
		Role: role,
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expirationMs) * time.Millisecond)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}

func (j *JWTUtil) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := JWTClaims{
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.refreshExpirationMs) * time.Millisecond)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}

func (j *JWTUtil) ParseToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

func (j *JWTUtil) ValidateToken(tokenStr string) bool {
	_, err := j.ParseToken(tokenStr)
	return err == nil
}

func (j *JWTUtil) ExpirationMs() int64        { return j.expirationMs }
func (j *JWTUtil) RefreshExpirationMs() int64 { return j.refreshExpirationMs }
