package jwt

import (
	"Foldr/models"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var JWTSecretKey = []byte("shayan")

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JWTSecretKey, nil
		})
		fmt.Print(token)
		if err != nil || !token.Valid {
			fmt.Print(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized-1"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Print(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized-2"})
			c.Abort()
			return
		}
		userID := claims["Username"].(string)
		// Need to check if the userID token is in the database by making a call to the database function?

		c.Set("userID", userID)

		c.Next()
	}
}
