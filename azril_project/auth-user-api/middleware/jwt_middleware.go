// middleware/jwt_middleaware.go
package middleware

import (
    "auth-user-api/controllers"
    "auth-user-api/domains"
    "auth-user-api/services"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v4"
    "github.com/labstack/echo/v4"
)

type JWTMiddlewareConfig struct {
    UserService services.UserService // Inject UserService
}

func NewJWTMiddleware(userService services.UserService) *JWTMiddlewareConfig {
    return &JWTMiddlewareConfig{UserService: userService}
}

func (mw *JWTMiddlewareConfig) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(ctx echo.Context) error {
        tokenString := ctx.Request().Header.Get("Authorization")

        if tokenString == "" {
            response := domains.BaseResponse{
                Code:    "401",
                Message: "Missing Authorization header",
                Error:   "No Authorization header provided",
            }
            return ctx.JSON(http.StatusUnauthorized, response)
        }

        if !strings.HasPrefix(tokenString, "Bearer ") {
            response := domains.BaseResponse{
                Code:    "401",
                Message: "Invalid Authorization header format",
                Error:   "Bearer token format error",
            }
            return ctx.JSON(http.StatusUnauthorized, response)
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        claims := &controllers.JWTClaims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte("my_secret_key"), nil
        })

        if err != nil || !token.Valid {
            response := domains.BaseResponse{
                Code:    "401",
                Message: "Invalid token",
                Error:   "Token validation error",
            }
            return ctx.JSON(http.StatusUnauthorized, response)
        }

        // Cek apakah user ada di database
        user, err := mw.UserService.GetUserByUsername(claims.Username)
        if err != nil || user == nil {
            response := domains.BaseResponse{
                Code:    "401",
                Message: "Invalid token - user not found",
                Error:   "User not found or deleted",
            }
            return ctx.JSON(http.StatusUnauthorized, response)
        }

        // Set username ke dalam context jika valid
        ctx.Set("username", claims.Username)

        // Lanjutkan ke handler berikutnya
        return next(ctx)
    }
}
