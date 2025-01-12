package util

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

type tokenCreateParams struct {
	Username string
	Email    string
	UserID   string
	jwt.StandardClaims
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorisation header provided")})
			c.Abort()
			return
		}

		claims, err := ValidateToken(clientToken)
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			c.Abort()
			return
		}

		c.Set("userid", claims.UserID)
		c.Next()
	}
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
		return nil, err
	}

	claims, ok := token.Claims.(*tokenCreateParams)

	if !ok {
		return nil, errors.New("the token is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}
