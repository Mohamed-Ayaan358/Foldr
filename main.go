package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// User struct to represent user data
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Sample database to store user data (Replace this with your database implementation)
var users = []User{
	{ID: 1, Username: "user1", Password: "password1"},
	{ID: 2, Username: "user2", Password: "password2"},
}

// JWT secret key
var jwtKey = []byte("secret_key")

// JWTClaims struct for JWT claims
type JWTClaims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token
func GenerateJWT(user User) (string, error) {
	claims := &JWTClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Middleware function to authenticate requests using JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	// User registration endpoint
	r.POST("/register", func(c *gin.Context) {
		var newUser User
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Store user in the database (Replace this with your database implementation)

		// Generate JWT token
		token, err := GenerateJWT(newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		users = append(users, newUser)
		fmt.Print(users)

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	// Example protected endpoint
	r.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Protected endpoint accessed"})
	})

	r.Run(":8080")
}
