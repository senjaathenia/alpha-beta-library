// controllers/user_controller.go

package controllers

import (
    "net/http"
    "time"
    "github.com/golang-jwt/jwt/v4"
    "auth-user-api/services"
    "auth-user-api/domains"
    "github.com/labstack/echo/v4"
)

type UserController struct {
    service services.UserService
}

func NewUserController(service services.UserService) *UserController {
    return &UserController{service}
}

// Register User godoc
func (c *UserController) RegisterUser(ctx echo.Context) error {
    type RegisterRequest struct {
        Username  string `json:"username" validate:"required"`
        Email     string `json:"email" validate:"required,email"`
        Password1 string `json:"password_1" validate:"required"`
        Password2 string `json:"password_2" validate:"required"`
    }

    var req RegisterRequest
    if err := ctx.Bind(&req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Failed processing input, try again. Error: " + err.Error(),
            Error:   "Binding error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if req.Username == "" {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Username cannot be empty. Field: username",
            Error:   "Validation error",
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if req.Email == "" {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Email cannot be empty. Field: email",
            Error:   "Validation error",
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if req.Password1 == "" {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Password 1 cannot be empty. Field: password_1",
            Error:   "Validation error",
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if req.Password2 == "" {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Password 2 cannot be empty. Field: password_2",
            Error:   "Validation error",
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := ctx.Validate(req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Validation error. Field: " + err.Error(),
            Error:   "Validation error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := c.service.Register(req.Username, req.Email, req.Password1, req.Password2); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Registration failed. Error: " + err.Error(),
            Error:   "Service error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    userResponse := domains.RegisterResponse{
        Username:  req.Username,
        Email:     req.Email,
        Password1: req.Password1,
        Password2: req.Password2,
    }

    response := domains.BaseResponse{
        Code:      "200",
        Message:   "User successfully registered",
        Data:      userResponse,
        Parameter: "username", 
    }    
    return ctx.JSON(http.StatusOK, response)
}

// Get All Users godoc
func (c *UserController) GetAllUsers(ctx echo.Context) error {
    users, err := c.service.GetAllUsers()
    if err != nil {
        response := domains.BaseResponse{
            Code:    "500",
            Message: "Failed to retrieve users. Error: " + err.Error(),
            Error:   "Service error: " + err.Error(),
        }
        return ctx.JSON(http.StatusInternalServerError, response)
    }

    response := domains.BaseResponse{
        Code:    "200",
        Message: "Users retrieved successfully",
        Data:    users,
    }
    return ctx.JSON(http.StatusOK, response)
}

// Update User godoc
func (c *UserController) UpdateUser(ctx echo.Context) error {
    type UpdateRequest struct {
        Username  string `json:"username" validate:"required"`
        Email     string `json:"email"`
        Password1 string `json:"password_1"`
        Password2 string `json:"password_2"`
    }

    userID := ctx.Param("id")
    if userID == "" {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "User ID is required. Field: id",
            Error:   "Validation error",
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    existingUser, err := c.service.GetUserByID(userID)
    if err != nil {
        response := domains.BaseResponse{
            Code:    "404",
            Message: "User not found. UserID: " + userID,
            Error:   "User retrieval error: " + err.Error(),
        }
        return ctx.JSON(http.StatusNotFound, response)
    }

    var req UpdateRequest
    if err := ctx.Bind(&req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Failed processing input. Error: " + err.Error(),
            Error:   "Binding error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    err = c.service.Update(userID, req.Username, req.Email, req.Password1, req.Password2)
    if err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Failed to update user. Error: " + err.Error(),
            Error:   "Service error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    userResponse := domains.UserResponse{
        UserID:   existingUser.ID,
        Username: req.Username,
        Email:    req.Email,
        Password: existingUser.Password,
    }

    response := domains.BaseResponse{
        Code:      "200",
        Message:   "User successfully updated. UserID: " + userID,
        Data:      userResponse,
        Parameter: "username", 
    }    
    return ctx.JSON(http.StatusOK, response)
}

// Delete User godoc
func (c *UserController) DeleteUser(ctx echo.Context) error {
    type DeleteRequest struct {
        UserID string `json:"user_id" validate:"required"`
    }

    var req DeleteRequest
    if err := ctx.Bind(&req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Delete failed. Error: " + err.Error(),
            Error:   "Binding error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := ctx.Validate(req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Delete failed. Validation error: " + err.Error(),
            Error:   "Validation error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    user, err := c.service.GetUserByID(req.UserID)
    if err != nil {
        response := domains.BaseResponse{
            Code:    "404",
            Message: "User not found. UserID: " + req.UserID,
            Error:   "User retrieval error: " + err.Error(),
        }
        return ctx.JSON(http.StatusNotFound, response)
    }

    if user.DeletedAt.Valid {
        response := domains.BaseResponse{
            Code:    "404",
            Message: "User not found. UserID: " + req.UserID,
        }
        return ctx.JSON(http.StatusNotFound, response)
    }

    if err := c.service.Delete(req.UserID); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Failed to delete user. UserID: " + req.UserID + ", Error: " + err.Error(),
            Error:   "Service error: " + err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    response := domains.BaseResponse{
        Code:      "200",
        Message:   "User deleted successfully. UserID: " + req.UserID,
        Data:      domains.DeleteResponse{UserID: req.UserID},
        Parameter: "user_id", 
    }    
    return ctx.JSON(http.StatusOK, response)
}

var jwtKey = []byte("my_secret_key")  // Pastikan menggunakan secret key yang sama

type JWTClaims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

// Login User
func (c *UserController) LoginUser(ctx echo.Context) error {
    type LoginRequest struct {
        Username string `json:"username" validate:"required"`
        Password string `json:"password" validate:"required"`
    }

    var req LoginRequest
    if err := ctx.Bind(&req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Invalid input",
            Error:  err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    if err := ctx.Validate(req); err != nil {
        response := domains.BaseResponse{
            Code:    "400",
            Message: "Validation error",
            Error:  err.Error(),
        }
        return ctx.JSON(http.StatusBadRequest, response)
    }

    // Authenticate the user
    err := c.service.Authenticate(req.Username, req.Password)
    if err != nil {
        if err.Error() == "user not found" {
            response := domains.BaseResponse{
                Code:    "404",
                Message: "User not found.",
                Error:   "UserNotFoundError",
            }
            return ctx.JSON(http.StatusNotFound, response)
        } else if err.Error() == "invalid username or password" {
            response := domains.BaseResponse{
                Code:    "401",
                Message: "Invalid username or password",
                Error:   "AuthenticationError",
            }
            return ctx.JSON(http.StatusUnauthorized, response)
        }
    
        response := domains.BaseResponse{
            Code:    "500",
            Message: "Internal server error",
            Error:   err.Error(),
        }
        return ctx.JSON(http.StatusInternalServerError, response)
    }    

    // Membuat token JWT dengan durasi 24 jam
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &JWTClaims{
        Username: req.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        response := domains.BaseResponse{
            Code:    "500",
            Message: "Failed to generate token",
            Error:  err.Error(),
        }
        return ctx.JSON(http.StatusInternalServerError, response)
    }

    response := domains.BaseResponse{
        Code:    "200",
        Message: "Successful login",
        Data: map[string]interface{}{
            "token": tokenString,
        },
        Error: "",
    }
    
    // Panggil helper function untuk memformat error jika kosong
    response.FormatError()
    
    return ctx.JSON(http.StatusOK, response)
}

// Route yang diproteksi
func (c *UserController) HelloProtected(ctx echo.Context) error {
    username := ctx.Get("username")
    if username == nil {
        response := map[string]string{
            "Message": "Unauthorized access. Missing or invalid token.",
        }
        return ctx.JSON(http.StatusUnauthorized, response)
    }

    response := map[string]string{
        "Message": "Hello, you have accessed a protected route!",
    }

    return ctx.JSON(http.StatusOK, response)
}
