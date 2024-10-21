// controllers/user_controller.go

package controllers

import (
    "net/http"
    "auth-user-api/services"

    "github.com/labstack/echo/v4"
)

type UserController struct {
    service services.UserService
}

func NewUserController(service services.UserService) *UserController {
    return &UserController{service}
}

// RegisterUser godoc
func (c *UserController) RegisterUser(ctx echo.Context) error {
    type RegisterRequest struct {
        Username  string `json:"username" validate:"required"`
        Email     string `json:"email" validate:"required,email"`
        Password1 string `json:"password_1" validate:"required"`
        Password2 string `json:"password_2" validate:"required"`
    }

    var req RegisterRequest
    if err := ctx.Bind(&req); err != nil {
        return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Failed processing input, try again"})
    }

    // Validasi input
    if err := ctx.Validate(req); err != nil {
        if _, ok := err.(*echo.HTTPError); ok {
            if req.Username == "" {
                return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Username cannot be empty"})
            }
            if req.Email == "" {
                return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Email cannot be empty"})
            }
            if req.Password1 == "" {
                return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Password 1 cannot be empty"})
            }
            if req.Password2 == "" {
                return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Password 2 cannot be empty"})
            }
        }
        return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    if err := c.service.Register(req.Username, req.Email, req.Password1, req.Password2); err != nil {
        return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    // Mengembalikan respon sukses dengan format yang diinginkan
    response := echo.Map{
        "message": "User successfully registered.",
        "user": echo.Map{
            "username":   req.Username,
            "email":      req.Email,
            "password_1": req.Password1,
            "password_2": req.Password2,
        },
        "code": "200",
    }

    return ctx.JSON(http.StatusOK, response)
}

// GetAllUsers godoc
func (c *UserController) GetAllUsers(ctx echo.Context) error {
    users, err := c.service.GetAllUsers()
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve users"})
    }

    // Jika tidak ada error, kembalikan semua data user dalam bentuk JSON
    return ctx.JSON(http.StatusOK, echo.Map{
        "message": "Users retrieved successfully",
        "users":   users,
        "code": "200",
    })
}

// UpdateUser mengupdate informasi pengguna
func (c *UserController) UpdateUser(ctx echo.Context) error {
    type UpdateRequest struct {
        Username  string `json:"username" validate:"required"`
        Email     string `json:"email"`
        Password1 string `json:"password_1"`
        Password2 string `json:"password_2"`
    }

    // Get the user ID from the URL parameter
    userID := ctx.Param("id")
    if userID == "" {
        return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is required"})
    }

    // Check if the user exists
    existingUser, err := c.service.GetUserByID(userID)
    if err != nil {
        return ctx.JSON(http.StatusNotFound, echo.Map{
            "error": "User not found",
            "code": "404",
        })
    }

    var req UpdateRequest
    if err := ctx.Bind(&req); err != nil {
        return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Failed processing input, try again"})
    }

    // Validasi input
    if err := ctx.Validate(req); err != nil {
        return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    // Update user di database
    err = c.service.Update(userID, req.Username, req.Email, req.Password1, req.Password2)
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    // Mengembalikan respons dengan format yang diminta
    return ctx.JSON(http.StatusOK, echo.Map{
        "message": "User successfully updated.",
        "user": echo.Map{
            "user_id":   existingUser.ID,
            "username":  req.Username,
            "email":     req.Email,
            "password":  existingUser.Password, // diasumsikan password sudah di-hash di dalam database
        },
        "code": "200",
    })
}

// DeleteUser godoc
func (c *UserController) DeleteUser(ctx echo.Context) error {
    type DeleteRequest struct {
        UserID string `json:"user_id" validate:"required"`
    }

    var req DeleteRequest
    if err := ctx.Bind(&req); err != nil {
        return ctx.JSON(http.StatusBadRequest, echo.Map{
            "message":         "Delete failed",
            "error":           err.Error(),
            "code": "400",
        })
    }

    if err := ctx.Validate(req); err != nil {
        return ctx.JSON(http.StatusBadRequest, echo.Map{
            "message":         "Delete failed",
            "error":           err.Error(),
            "code": "400",
        })
    }

    // Check if the user exists and if the user has already been deleted
    user, err := c.service.GetUserByID(req.UserID)
    if err != nil {
        return ctx.JSON(http.StatusNotFound, echo.Map{
            "message":         "Delete failed",
            "error":           echo.Map{"error": "User not found"},
            "code": "404",
        })
    }

    // Cek apakah user sudah dihapus sebelumnya
    if user.DeletedAt.Valid {
        return ctx.JSON(http.StatusNotFound, echo.Map{
            "message":         "Delete failed",
            "error":           echo.Map{"error": "User not found"},
            "code": "404",
        })
    }

    // Delete the user
    if err := c.service.Delete(req.UserID); err != nil {
        return ctx.JSON(http.StatusBadRequest, echo.Map{
            "message":         "Delete failed",
            "error":           err.Error(),
            "code": "400",
        })
    }

    return ctx.JSON(http.StatusOK, echo.Map{
        "message":         "User deleted successfully",
        "code": "200",
    })
}

// LoginUser godoc
func (c *UserController) LoginUser(ctx echo.Context) error {
    type LoginRequest struct {
        Username string `json:"username" validate:"required"`
        Password string `json:"password" validate:"required"`
    }

    var req LoginRequest
    if err := ctx.Bind(&req); err != nil {
        return ctx.JSON(http.StatusBadRequest, echo.Map{
            "message":         "Authenticate failed",
            "error":           err.Error(),
            "code": "400",
        })
    }

    if err := ctx.Validate(req); err != nil {
        return ctx.JSON(http.StatusBadRequest, echo.Map{
            "message":         "Authenticate failed",
            "error":           err.Error(),
            "code": "400",
        })
    }

    // Authenticate user
    err := c.service.Authenticate(req.Username, req.Password)
    if err != nil {
        if err.Error() == "user not found" {
            return ctx.JSON(http.StatusNotFound, echo.Map{
                "message":         "Authentication failed",
                "errortype":           echo.Map{"error-type": "User not found"},
                "code": "404",
            })
        }
        return ctx.JSON(http.StatusUnauthorized, echo.Map{
            "message":         "Authentication failed",
            "error":           echo.Map{"error-type": "Invalid Username or Password"},
            "code": "401",
        })
    }

    return ctx.JSON(http.StatusOK, echo.Map{
        "message":         "Authentication success",
        "code": "200",
    })
}
