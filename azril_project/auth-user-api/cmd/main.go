package main

import (
    "fmt"
    "log"
    "auth-user-api/controllers"
    "auth-user-api/repository"
    "auth-user-api/services"
    "auth-user-api/models"
    "auth-user-api/utils"
    "auth-user-api/middleware"  // Tambahkan ini

    "github.com/labstack/echo/v4"
    echoMiddleware "github.com/labstack/echo/v4/middleware"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    // Konfigurasi Database
    dsn := "host=localhost user=postgres password=arnoarno dbname=api-auth port=5432 sslmode=disable TimeZone=Asia/Jakarta"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Jalankan Migrasi
    err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";").Error
    if err != nil {
        log.Fatalf("Failed to create extension: %v", err)
    }

    err = db.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    // Inisialisasi Repository, Service, dan Controller
    userRepo := repository.NewUserRepository(db)
    userService := services.NewUserService(userRepo)
    userController := controllers.NewUserController(userService)

    // Inisialisasi Echo
    e := echo.New()

    // Middleware
    e.Use(echoMiddleware.Logger())
    e.Use(echoMiddleware.Recover())

    // Validator
    e.Validator = utils.NewValidator()

    // Routes
    e.POST("/register", userController.RegisterUser)
    e.POST("/login", userController.LoginUser)
    e.GET("/users", userController.GetAllUsers)
    e.PUT("/update/:id", userController.UpdateUser)
    e.DELETE("/delete", userController.DeleteUser)
    
    jwtMiddleware := middleware.NewJWTMiddleware(userService)

    // Rute dengan middleware JWT
    e.GET("/protected/hello", userController.HelloProtected, jwtMiddleware.JWTMiddleware)

    // Start Server
    port := "8080"
    fmt.Printf("Server running on port %s\n", port)
    if err := e.Start(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
