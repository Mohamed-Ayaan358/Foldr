package database

import (
	"Foldr/models"
	"database/sql"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "ayaan2002"
	dbname   = "foldr"
)

var db *sql.DB

func getDatabaseConnection() (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", "host="+host+" port="+strconv.Itoa(port)+" user="+user+" password="+password+" dbname="+dbname+" sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	return db, err
}

func Databaseregister(newUser models.User, hashedPassword []byte) error {
	var err error
	db, _ = getDatabaseConnection()
	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", newUser.Username, string(hashedPassword))
	return err
}

func Databaselogin(loginUser models.User, dbUser models.User) (error, models.User) {
	var err error
	db, _ = getDatabaseConnection()
	err = db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", loginUser.Username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)

	return err, dbUser
}
