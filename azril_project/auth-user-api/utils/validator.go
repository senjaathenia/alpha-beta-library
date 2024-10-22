// utils/validator.go

package utils

import (
    "errors"
    "regexp"
	"github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
    "net/http"
)

// CustomValidator adalah implementasi dari echo.Validator
type CustomValidator struct {
    validator *validator.Validate
}

// NewValidator mengembalikan instance baru dari CustomValidator
func NewValidator() *CustomValidator {
    return &CustomValidator{validator: validator.New()}
}

// Validate melakukan validasi terhadap struct menggunakan go-playground/validator
func (cv *CustomValidator) Validate(i interface{}) error {
    if err := cv.validator.Struct(i); err != nil {
        // Mengonversi error validasi ke HTTP Error
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    return nil
}

// isAlphabetic checks if the given string contains only alphabetic characters
func isAlphabetic(s string) bool {
    for _, char := range s {
        if !('a' <= char && char <= 'z' || 'A' <= char && char <= 'Z') {
            return false
        }
    }
    return true
}

func ValidatePassword(password string) error {
    var (
        minLength  = 8
        hasUpper   = regexp.MustCompile(`[A-Z]`)
        hasNumber  = regexp.MustCompile(`[0-9]`)
        hasSpecial = regexp.MustCompile(`[!@#~$%^&*()+|_]{1}`)
    )

    if len(password) < minLength {
        return errors.New("password harus minimal 8 karakter")
    }
    if !hasUpper.MatchString(password) {
        return errors.New("password harus mengandung minimal 1 huruf besar")
    }
    if !hasNumber.MatchString(password) {
        return errors.New("password harus mengandung minimal 1 angka")
    }
    if !hasSpecial.MatchString(password) {
        return errors.New("password harus mengandung minimal 1 simbol")
    }

    // Alfanumerik + simbol sudah dipenuhi dengan pengecekan di atas
    return nil
}
