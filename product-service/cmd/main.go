package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ajitirto/ms/product-service/internal/infrastructure/postgres"
	"github.com/ajitirto/ms/product-service/internal/usecase"
	"github.com/ajitirto/ms/product-service/pkg/server"


	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file (defaulting to existing environment vars)")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" 
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	
	if err != nil {
		log.Fatal(fmt.Errorf("error opening database: %w", err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(fmt.Errorf("error connecting to database: %w", err))
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	// 3. Dependency Injection
	// Infrastructure -> Usecase -> Delivery
	productRepo := postgres.NewProductRepository(db)
	productUC := usecase.NewProductService(productRepo)
	productHandler := server.NewProductHandler(productUC)

	http.HandleFunc("/", productHandler.StatusHandler)
	http.HandleFunc("/products", productHandler.HandleProducts)  
	http.HandleFunc("/products/", productHandler.GetOrderByID) 

	log.Printf("Server starting on :%s\n", appPort)

	if err := http.ListenAndServe(":"+appPort, nil); err != nil {
		log.Fatal(err)
	}
}
