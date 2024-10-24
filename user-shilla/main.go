package main

import (
	"log"
	"project-golang-crud/domains"
	"project-golang-crud/pkg/config"
	"project-golang-crud/pkg/delivery"
	"project-golang-crud/pkg/repository"
	"project-golang-crud/pkg/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func main() {
	db := config.ConnectDB()
	if db == nil {
		log.Fatal("Database connection failed")
	}
	log.Println("Database connection successfully")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	migrate(db)

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	delivery.NewUserHandler(e,userUsecase)

	e.Logger.Fatal(e.Start(":8082"))
}

func migrate(db *gorm.DB)  {
	err := db.AutoMigrate(&domains.User{})
	if err != nil {
		log.Fatalf("Error in database migration: %v", err)
	}
	log.Println("Database migration completed!")
}