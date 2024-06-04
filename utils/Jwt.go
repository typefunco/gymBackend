package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func jwtSecret() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	secretKey := os.Getenv("secretKey")
	return secretKey
}

func GenerateToken(UseriId int, username string, isSuperUser bool) (string, error) {
	secretKey := jwtSecret()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":        username,
		"id":           UseriId,
		"is_superuser": isSuperUser,
		"exp":          time.Now().Add(time.Hour * 3).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int, error) {
	secretKey := jwtSecret()
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("token not valid")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, errors.New("COULD'T PARSE TOKEN")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return 0, errors.New("INVALID TOKEN")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("INVALID CLAIMS")
	}

	// login := claims["login"].(string)
	userId := int(claims["id"].(float64)) // actual data of type float64 and I convert data to int

	return userId, nil
}
