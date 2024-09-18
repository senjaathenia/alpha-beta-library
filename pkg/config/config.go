	package config

	import (
		"log"
		"os"

		"github.com/joho/godotenv"
		"gorm.io/driver/postgres"
		"gorm.io/gorm"
	)

	// ConnectDB initializes and returns a PostgreSQL connection using GORM
	func ConnectDB() *gorm.DB {
		// Load environment variables from .env file in the conf folder
		err := godotenv.Load("conf/config.env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		// Get the database URL from environment variable
		dsn := os.Getenv("DB_URL")
		if dsn == "" {
			log.Fatal("DB_URL is not set in the environment")
		}

		// Connect to PostgreSQL database using GORM
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}

		log.Println("Database connection successfully established!")
		return db
	}
