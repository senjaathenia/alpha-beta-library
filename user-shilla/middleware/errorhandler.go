package middleware

import (
	"net/http"
	"project-golang-crud/domains"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			// Memeriksa jika endpoint tidak ditemukan
			if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusNotFound {
				return c.JSON(http.StatusBadRequest, domains.Response{
					Message: "Endpoint tidak ditemukan",
					Data:    nil,
					Errors: []domains.ErrorDetail{
						{Message: "Silakan periksa URL yang diminta", Parameter: "url"},
					},
					Code: http.StatusBadRequest,
				})
			}
			// Memeriksa kesalahan lainnya
			return err
		}
		return nil
	}
}
