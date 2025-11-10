package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file (defaulting to existing environment vars)")
    }
    //
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	appPort := os.Getenv("APP_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening database: %w", err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to database: %w", err))
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Product Service is running and connected to DB!")
	})

	log.Printf("Server starting on :%s\n", appPort)
	if err := http.ListenAndServe(":"+appPort, nil); err != nil {
		log.Fatal(err)
	}
}