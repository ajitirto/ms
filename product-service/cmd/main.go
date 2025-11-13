package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/ajitirto/ms/product-service/internal/infrastructure/postgres"
	redisInfra "github.com/ajitirto/ms/product-service/internal/infrastructure/redis"
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

	redisHost := os.Getenv("REDIS_HOST") 
	redisPort := os.Getenv("REDIS_PORT") 

	// fmt.Println(redisHost, redisPort)


	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		// Addr:     "127.0.0.1:6379",
		Password: "", // tidak ada password
		DB:       0,  // gunakan DB default
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Could not connect to Redis at %s:%s. Running without cache. Error: %v", redisHost, redisPort, err)
	} else {
		log.Println("Successfully connected to Redis cache.")
	}

	const cacheTTL = 5 * time.Minute 
	productCache := redisInfra.NewProductCacheRedis(redisClient, cacheTTL) 


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

	

	productRepo := postgres.NewProductRepository(db)
	productUC := usecase.NewProductService(productRepo, productCache)
	productHandler := server.NewProductHandler(productUC)

	http.HandleFunc("/", productHandler.StatusHandler)
	http.HandleFunc("/products", productHandler.HandleProducts)  
	http.HandleFunc("/products/", productHandler.GetOrderByID) 

	log.Printf("Server starting on :%s\n", appPort)

	if err := http.ListenAndServe(":"+appPort, nil); err != nil {
		log.Fatal(err)
	}
}
