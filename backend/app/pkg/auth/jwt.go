package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// type HashType int64

// const (
// 	ACTIVATE HashType = iota
// 	PASSWORD_RESET
// )

type AuthJwtData struct {
	Email string `json:"email,omitempty"`
	Id    uint16 `json:"id,omitempty"`
}

type LinkJwtData struct {
	Id uint16 `json:"id,omitempty"`
	// Type HashType
}

type JwtData interface {
	AuthJwtData | LinkJwtData
}

type Jwt[T any] struct {
	jwt.StandardClaims
	Data T
}

type AuthJwt = Jwt[AuthJwtData]
type LinkJwt = Jwt[LinkJwtData]

func Encode[T JwtData](claims *Jwt[T], expireMins int) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Minute * time.Duration(expireMins)).Unix()
	claims.StandardClaims.Issuer = "smer-auth"

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("encodeJwt: %v", err)
	}
	return tokenString, nil
}

func Decode[T JwtData](claims *Jwt[T], s string) (*jwt.Token, *Jwt[T], error) {
	token, err := jwt.ParseWithClaims(s, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("decodeJwt: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		if err == nil {
			err = ErrInvalidToken
		}

		if strings.Contains(err.Error(), "expired") {
			return nil, nil, ErrExpiredToken
		}

		return nil, nil, fmt.Errorf("decodeJwt: %v", err)
	}
	claims, ok := token.Claims.(*Jwt[T])
	if !ok {
		return nil, nil, errors.New("decodeJwt: invalid claims")
	}
	return token, claims, nil
}
