package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codepnw/go-car-management/database"
	"github.com/codepnw/go-car-management/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	envFile = "dev.env"
	version = "/v1"
	schemaFile = "database/schema.sql"
)

func main() {
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	// Database
	db := database.ConnectPostgres(os.Getenv("DB_CONN_STR"))
	defer db.Close()

	if err := database.ExecuteSQLSchema(db, schemaFile); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// Routes
	routes.NewRoutes(db, r, version)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("server starting on port:", port)
	r.Run(":" + port)
}
