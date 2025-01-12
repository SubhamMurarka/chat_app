package Middleware

import (
	"fmt"
	"net/http"
	"strings"

	util "github.com/SubhamMurarka/chat_app/User/Util"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorisation header provided")})
			c.Abort()
			return
		}

		claims, err := util.ValidateToken(clientToken)
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("userid", claims.UserID)
		c.Next()
	}
}
