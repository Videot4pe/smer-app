package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type BaseJwt struct {
	jwt.StandardClaims
}

type AuthJwt struct {
	BaseJwt
	Email string `json:"email,omitempty"`
	Id    uint16 `json:"id,omitempty"`
}

type HashType int64

const (
	ACTIVATE HashType = iota
	PASSWORD_RESET
)

type LinkJwt struct {
	BaseJwt
	id   uint64
	Type HashType
}

type JwtInterface[T BaseJwt] interface {
	Encode(expireMins int) (string, error)
	Decode(jwtString string) (*jwt.Token, *T, error)
}

func (claims *BaseJwt) Encode(expireMins int) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Minute * time.Duration(expireMins)).Unix()
	claims.StandardClaims.Issuer = "smer-auth"

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("encodeJwt: %v", err)
	}
	return tokenString, nil
}

func Decode[T *BaseJwt](s string) (*jwt.Token, T, error) {
	token, err := jwt.ParseWithClaims(s, &BaseJwt{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("decodeJwt: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		if err == nil {
			err = errors.New("invalid token")
		}
		return nil, nil, fmt.Errorf("decodeJwt: %v", err)
	}
	claims, ok := token.Claims.(*BaseJwt)
	if !ok {
		return nil, nil, errors.New("decodeJwt: invalid claims")

	}
	return token, claims, nil
}
