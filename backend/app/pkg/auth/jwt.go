package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtClaims struct {
	jwt.StandardClaims
	Email string `json:"email,omitempty"`
	Id    uint16 `json:"id,omitempty"`
}

func NewJwtClaims(email string, id uint16) *JwtClaims {

	expireToken := time.Now().Add(time.Minute * 10).Unix()

	return &JwtClaims{
		Email: email,
		Id:    id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "smer-auth",
		},
	}
}

func (claims *JwtClaims) EncodeJwt() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("encodeJwt: %v", err)
	}
	return tokenString, nil
}

func DecodeJwt(s string) (*jwt.Token, *JwtClaims, error) {
	token, err := jwt.ParseWithClaims(s, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
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
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, nil, errors.New("decodeJwt: invalid claims")

	}
	return token, claims, nil
}
