package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"go_todolist/internal/entity"
)

type Claims struct {
	IDUsers int64  `json:"id_users"`
	Nama    string `json:"nama"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(secret string, expiresIn time.Duration, user entity.User) (string, error) {
	now := time.Now()
	claims := Claims{
		IDUsers: user.IDUsers,
		Nama:    user.Nama,
		Email:   user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go_todolist",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(secret, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing tidak valid")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token tidak valid")
	}
	return claims, nil
}
