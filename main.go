package main

import (
	"Foldr/jwt"
	"Foldr/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "ayaan2002"
	dbname   = "foldr"
)

var db *sql.DB

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	var err error
	db, err = sql.Open("postgres", "host="+host+" port="+strconv.Itoa(port)+" user="+user+" password="+password+" dbname="+dbname+" sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()
	r := gin.Default()

	r.POST("/register", registerUser)
	r.POST("/login", loginUser)
	r.GET("/protected", jwt.AuthMiddleware(), protectedEndpoint)

	r.Run(":8080")
	fmt.Print("Server being hosted on 8080")
}

// Handler function to register a new user
func registerUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

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
	var loginUser models.User
	if err := c.BindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var dbUser models.User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", loginUser.Username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginUser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := jwt.GenerateJWT(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func protectedEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Protected endpoint accessed"})
}
