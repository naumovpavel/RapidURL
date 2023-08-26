package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const expirationInterval = time.Hour * 24 * 30

var secretKey, _ = os.LookupEnv("JWT_SECRET_KEY")

func CreateJWT(id int) (string, error) {
	const op = "auth.CreateJWT"
	secretKey := []byte(secretKey)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userid"] = id
	claims["exp"] = time.Now().Add(expirationInterval).Unix() // Expiration time

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error creating token:", err)
		return "", fmt.Errorf("%s : %w", op, err)
	}

	return signedToken, nil
}

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

func DecodeJWT(jwtToken string) (int, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := int(claims["userid"].(float64))
		expiration := int64(claims["exp"].(float64))
		if expiration > time.Now().Unix() {
			return userId, nil
		} else {
			return 0, ErrTokenExpired
		}
	} else {
		return 0, ErrTokenInvalid
	}
}
