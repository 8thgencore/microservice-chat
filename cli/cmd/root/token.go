package root

import (
	"errors"
	"io"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/8thgencore/microservice-common/pkg/closer"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

func readToken() (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	closer.Add(file.Close)

	token, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func getTokenClaims(token string) (*UserClaims, error) {
	t, err := jwt.ParseWithClaims(
		token,
		&UserClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return []byte(""), nil
		})
	if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return nil, err
	}

	claims, ok := t.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

func isTokenExpired(claims *UserClaims) error {
	expire, err := claims.RegisteredClaims.GetExpirationTime()
	if err != nil {
		return err
	}

	if expire.Before(time.Now()) {
		// Token expired
		return errors.New("access token expired")
	}

	return nil
}
