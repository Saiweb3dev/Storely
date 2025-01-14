package utils

import (
    "time"

    "github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your_secret_key")

func GenerateJWT(userID string) (string, error) {
    claims := jwt.MapClaims{
        "userID": userID,
        "exp":    time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
}
