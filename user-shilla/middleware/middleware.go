package middleware

import (
	"net/http"
	"os"
	"project-golang-crud/domains"
	"strings"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware function
func JWTMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Mendapatkan token dari header Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, domains.Response{
					Message: "Authorization header missing",
					Errors: []domains.ErrorDetail{
						{
							Message:  "Authorization header is required",
							Parameter: "Authorization",
						},
					},
					Code: http.StatusUnauthorized,
				})
			}

			// Memeriksa apakah token menggunakan skema Bearer
			tokenString := strings.Split(authHeader, " ")
			if len(tokenString) != 2 || tokenString[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, domains.Response{
					Message: "Invalid token format",
					Errors: []domains.ErrorDetail{
						{
							Message:  "Token must be provided in Bearer <token> format",
							Parameter: "Authorization",
						},
					},
					Code: http.StatusUnauthorized,
				})
			}

			// Parsing token JWT
			secretKey := []byte(os.Getenv("JWT_SECRET"))
			token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
				// Pastikan algoritma enkripsi token cocok dengan yang diharapkan
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
				}
				return secretKey, nil
			})

			// Jika token tidak valid
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, domains.Response{
					Message: "Invalid or expired token",
					Errors: []domains.ErrorDetail{
						{
							Message:  "The token provided is either invalid or has expired. Please login again to obtain a new token.",
							Parameter: "Authorization",
						},
					},
					Code: http.StatusUnauthorized,
				})
			}

			// Mengambil klaim dari token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return c.JSON(http.StatusUnauthorized, domains.Response{
					Message: "Invalid or expired token",
					Errors: []domains.ErrorDetail{
						{
							Message:  "The token provided is either invalid or has expired. Please login again to obtain a new token.",
							Parameter: "Authorization",
						},
					},
					Code: http.StatusUnauthorized,
				})
			}

			// Tambahkan logging untuk klaim
			fmt.Println("Claims:", claims)

			// Periksa role pengguna
			if len(allowedRoles) > 0 {
				role, ok := claims["role"].(string)
				if !ok {
					return c.JSON(http.StatusForbidden, domains.Response{
						Message: "Role claim missing or invalid",
						Errors: []domains.ErrorDetail{
							{
								Message:  "The role claim is either missing or not valid.",
								Parameter: "role",
							},
						},
						Code: http.StatusForbidden,
					})
				}

				allowed := false
				for _, allowedRole := range allowedRoles {
					if role == allowedRole {
						allowed = true
						break
					}
				}
				if !allowed {
					return c.JSON(http.StatusForbidden, domains.Response{
						Message: "Access forbidden: role required",
						Errors: []domains.ErrorDetail{
							{
								Message:  "You do not have permission to access this resource.",
								Parameter: "Authorization",
							},
						},
						Code: http.StatusForbidden,
					})
				}
			}

			// Token valid dan role sesuai, lanjutkan ke handler berikutnya
			return next(c)
		}
	}
}
