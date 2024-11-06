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

	bookRepo := repository.NewGenericRepository(db)
	bookUsecase := usecase.NewBookUsecase(bookRepo)
	delivery.NewBookHandler(e, bookUsecase)

	authorRepo := repository.NewGenericRepository(db)
	authorUsecase := usecase.NewAuthorUsecase(authorRepo)
	delivery.NewAuthorHandler(e, authorUsecase)

	publisherRepo := repository.NewGenericRepository(db)
	publisherUsecase := usecase.NewPublisherUsecase(publisherRepo)
	delivery.NewPublisherHandler(e, publisherUsecase)

	loansRepo := repository.NewGenericRepository(db)
	loansUsecase := usecase.NewLoanUsecase(loansRepo)
	delivery.NewLoanHandler(e, loansUsecase)

	userRepo := repository.NewUserRepository(db)

	loanRepo := repository.NewGenericRepository(db)
	loanUsecase := usecase.NewBookRequestUsecase(loanRepo)
	delivery.NewBookRequestHandler(e, loanUsecase)

	userUsecase := usecase.NewUserUsecase(userRepo, bookRepo, loanRepo, loanUsecase)
	delivery.NewUserHandler(e, userUsecase)

	e.Logger.Fatal(e.Start(":8082"))

}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&domains.User{}, &domains.Book{}, &domains.Author{}, &domains.Publisher{}, &domains.BookLoans{})
	if err != nil {
		log.Fatalf("Error in database migration: %v", err)
	}
	log.Println("Database migration completed!")
}
