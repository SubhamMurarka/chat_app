package Util

import (
	"log"
	"time"

	"github.com/SubhamMurarka/chat_app/User/Config"
	"github.com/golang-jwt/jwt"
)

type tokenCreateParams struct {
	Username string
	Email    string
	UserID   string
	jwt.StandardClaims
}

var SECRET_KEY string = Config.Config.JwtSecret

func GenerateAllTokens(userID string, username string, email string) (signedToken string, err error) {
	claims := &tokenCreateParams{
		Username: username,
		Email:    email,
		UserID:   userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Minute * time.Duration(10)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Printf("Error creating token %v", err)
		return "", ErrInternal
	}

	return token, nil
}

func ValidateToken(signedToken string) (*tokenCreateParams, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&tokenCreateParams{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		log.Printf("Error parsing token : %v", err)
		return nil, ErrInternal
	}

	claims, ok := token.Claims.(*tokenCreateParams)
	if !ok {
		log.Printf("Invalid token : %v", err)
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		log.Printf("Expired token : %v", err)
		return nil, ErrExpiredToken
	}

	return claims, nil
}
