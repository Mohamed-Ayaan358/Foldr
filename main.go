package main

import (
	"Foldr/routes"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
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

func main() {
	var err error
	db, err = sql.Open("postgres", "host="+host+" port="+strconv.Itoa(port)+" user="+user+" password="+password+" dbname="+dbname+" sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()
	r := gin.Default()
	// config := cors.DefaultConfig()
	// config.AllowCredentials = true
	// r.Use(cors.New(config))
	routes.RouteHandler(r)

	r.Run(":8080")
	fmt.Print("Server being hosted on 8080")
}
