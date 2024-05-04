package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// JWTSecretKey is a secret key used to sign JWT tokens
var JWTSecretKey = []byte("shayan")

// GenerateJWT generates a JWT token for a given user
func GenerateJWT(user User) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key and get the complete encoded token as a string
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Middleware function to validate JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT token from Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JWTSecretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "ayaan2002"
	dbname   = "foldr"
)

// Initialize a global database connection variable
var db *sql.DB

// User struct to represent user data
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// Create the database connection
	var err error
	db, err = sql.Open("postgres", "host="+host+" port="+strconv.Itoa(port)+" user="+user+" password="+password+" dbname="+dbname+" sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Initialize the Gin router
	r := gin.Default()

	// Add endpoints for user registration and login
	r.POST("/register", registerUser)
	r.POST("/login", loginUser)
	r.GET("/protected", protectedEndpoint)

	// Run the Gin server
	r.Run(":8080")
	fmt.Print("Server being hosted on 8080")
}

// Handler function to register a new user
func registerUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Insert the user into the database
	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", newUser.Username, string(hashedPassword))
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}
	fmt.Print("Signed up")

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Handler function to authenticate a user
func loginUser(c *gin.Context) {
	var loginUser User
	if err := c.BindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Query the database to retrieve the user with the given username
	var dbUser User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", loginUser.Username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the hashed password from the database with the password provided by the user
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginUser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Example of a protected endpoint that requires authentication
func protectedEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Protected endpoint accessed"})
}
