package delivery

import (
	"errors"
	"net/http"
	"project-golang-crud/domains"
	"strings"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type UserHandler struct {
	Usecase domains.UserUsecase
}

func NewUserHandler(e *echo.Echo, u domains.UserUsecase) {
	handler := &UserHandler{Usecase: u}

	e.POST("/register", handler.Register)
	e.PUT("/update/:id", handler.Update)
	e.DELETE("/delete", handler.Delete)
	e.POST("/validate", handler.Validate)
}

func (h *UserHandler) Register(c echo.Context) error {
    var req struct {
        Username  interface{} `json:"username"`
        Email     interface{} `json:"email"`
        Password1 interface{} `json:"password_1"`
        Password2 interface{} `json:"password_2"`
    }

    // Bind request body ke struct
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Message: "Invalid Request",
            Errors: []domains.ErrorDetail{
                {Message: "Failed to parse request body", Parameter: "Request Body"},
            },
            Code: http.StatusBadRequest,
        })
    }

    // Slice untuk menyimpan error
    var validationErrors []domains.ErrorDetail

    // Cek tipe data untuk setiap field
    if _, ok := req.Username.(string); !ok {
        validationErrors = append(validationErrors, domains.ErrorDetail{
            Message: "Field must be a string", 
            Parameter: "username",
        })
    }
    if _, ok := req.Email.(string); !ok {
        validationErrors = append(validationErrors, domains.ErrorDetail{
            Message: "Field must be a string", 
            Parameter: "email",
        })
    }
    if _, ok := req.Password1.(string); !ok {
        validationErrors = append(validationErrors, domains.ErrorDetail{
            Message: "Field must be a string", 
            Parameter: "password_1",
        })
    }
    if _, ok := req.Password2.(string); !ok {
        validationErrors = append(validationErrors, domains.ErrorDetail{
            Message: "Field must be a string", 
            Parameter: "password_2",
        })
    }

    // Jika ada error validasi, kembalikan respons dengan semua error
    if len(validationErrors) > 0 {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Message: "Validation Errors",
            Errors:   validationErrors,
            Code:    http.StatusBadRequest,
        })
    }

    // Cek apakah password1 dan password2 cocok
    if req.Password1 != req.Password2 {
        validationErrors = append(validationErrors, domains.ErrorDetail{
            Message: "Passwords don't match", 
            Parameter: "password",
        })
    }

    // Panggil usecase untuk registrasi
    user, err := h.Usecase.Register(req.Username.(string), req.Email.(string), req.Password1.(string))
    if err != nil {
        // Jika validasi gagal, tampilkan semua error validasi
        // Misalkan error dari `usecase` berisi beberapa error
        for _, msg := range strings.Split(err.Error(), "; ") {
            if strings.Contains(msg, "Username") {
                validationErrors = append(validationErrors, domains.ErrorDetail{
                    Message: msg,
                    Parameter: "username",
                })
            } else if strings.Contains(msg, "Invalid email") {
                validationErrors = append(validationErrors, domains.ErrorDetail{
                    Message: msg,
                    Parameter: "email",
                })
            } else if strings.Contains(msg, "Password") {
                validationErrors = append(validationErrors, domains.ErrorDetail{
                    Message: msg,
                    Parameter: "password",
                })
            }
        }
        
        return c.JSON(http.StatusBadRequest, domains.Response{
            Message: "Validation Errors",
            Errors: validationErrors,
            Code: http.StatusBadRequest,
        })
    }

    // Menyusun response dengan field deleted_at
    return c.JSON(http.StatusCreated, domains.Response{
        Message: "User created successfully",
        Data: domains.User{
            ID:        user.ID,
            Username:  user.Username,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
            DeletedAt: user.DeletedAt,
        },
        Code: http.StatusCreated,
    })
}

func (h *UserHandler) Update(c echo.Context) error {
    var req struct {
        Username  interface{} `json:"username"`  
        Email     interface{} `json:"email"`
        Password1 interface{} `json:"password_1"`
        Password2 interface{} `json:"password_2"`
    }

    id := c.Param("id")
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Message: "Invalid Request",
            Errors: []domains.ErrorDetail{
                {Message: "Failed to parse request body", Parameter: "Request Body"},
            },
            Code: http.StatusBadRequest,
        })
    }

    var validationErrors []domains.ErrorDetail

    // Validasi tipe data username harus string
    if _, ok := req.Username.(string); !ok {
        validationErrors = append(validationErrors, domains.ErrorDetail{
            Message: "Field must be a string",
            Parameter: "username",
        })
    }

    // Validasi opsional email
    var email string
    if req.Email != nil {
        if emailStr, ok := req.Email.(string); ok {
            email = emailStr
        } else {
            validationErrors = append(validationErrors, domains.ErrorDetail{
                Message: "Field must be a string",
                Parameter: "email",
            })
        }
    }

    // Validasi opsional password1 dan password2
    var password1, password2 string
    if req.Password1 != nil {
        if passwordStr, ok := req.Password1.(string); ok {
            password1 = passwordStr
        } else {
            validationErrors = append(validationErrors, domains.ErrorDetail{
                Message: "Field must be a string",
                Parameter: "password_1",
            })
        }
    }

    if req.Password2 != nil {
        if passwordStr, ok := req.Password2.(string); ok {
            password2 = passwordStr
        } else {
            validationErrors = append(validationErrors, domains.ErrorDetail{
                Message: "Field must be a string",
                Parameter: "password_2",
            })
        }
    }

    // Jika ada error validasi, kembalikan respons dengan semua error
    if len(validationErrors) > 0 {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Message: "Validation Errors",
            Errors:   validationErrors,
            Code:    http.StatusBadRequest,
        })
    }

    // Cek apakah password1 dan password2 cocok jika keduanya diisi
    if password1 != "" || password2 != "" {
        if password1 != password2 {
            validationErrors = append(validationErrors, domains.ErrorDetail{
                Message: "Passwords don't match",
                Parameter: "password",
            })
        }
    }

    // Panggil usecase untuk update
    err := h.Usecase.Update(id, req.Username.(string), email, password1)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.JSON(http.StatusNotFound, domains.Response{
                Message: "User not found",
                Code:    http.StatusNotFound,
            })
        }

        // Untuk kesalahan validasi yang berasal dari usecase
        for _, msg := range strings.Split(err.Error(), "; ") {
            if strings.Contains(msg, "Username") {
                validationErrors = append(validationErrors, domains.ErrorDetail{
                    Message: msg,
                    Parameter: "username",
                })
            } else if strings.Contains(msg, "Invalid email") {
                validationErrors = append(validationErrors, domains.ErrorDetail{
                    Message: msg,
                    Parameter: "email",
                })
            } else if strings.Contains(msg, "Password") {
                validationErrors = append(validationErrors, domains.ErrorDetail{
                    Message: msg,
                    Parameter: "password",
                })
            }
        }
        
        return c.JSON(http.StatusNotFound, domains.Response{
            Message: "User not found",
            Errors: validationErrors,
            Code: http.StatusNotFound,
        })
    }

    return c.JSON(http.StatusOK, domains.Response{
        Message: "User updated successfully",
        Code:    http.StatusOK,
    })
}

func (h *UserHandler) Delete(c echo.Context) error {
    var req domains.DeleteRequest

    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Message: "Invalid Request Body",
            Code:    http.StatusBadRequest,
        })
    }

    // Dapatkan pengguna yang dihapus
    user, err := h.Usecase.Delete(req.ID)
    if err != nil {
        if errors.Is(err, domains.ErrUserNotFound) {
            return c.JSON(http.StatusNotFound, domains.Response{
                Message: "User not found",
                Errors: nil,
                Code:    http.StatusNotFound,
            })
        }

        return c.JSON(http.StatusNotFound, domains.Response{
            Message: "User not found",
            Errors: nil,
            Code:    http.StatusNotFound,
        })
    }

    // Ambil kembali data user yang sudah dihapus untuk response
    updatedUser, err := h.Usecase.GetByID(user.ID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Message: "Failed to retrieve deleted user data",
            Code:    http.StatusInternalServerError,
        })
    }

    return c.JSON(http.StatusOK, domains.Response{
        Message: "User Deleted",
        Data: domains.User{
            ID:        updatedUser.ID,
            Username:  updatedUser.Username,
            Email:     updatedUser.Email,
            CreatedAt: updatedUser.CreatedAt,
            UpdatedAt: updatedUser.UpdatedAt,
            DeletedAt: updatedUser.DeletedAt, 
        },
        Errors: nil,
        Code: http.StatusOK,
    })
}

func (h *UserHandler) Validate(c echo.Context) error {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    // Bind JSON request ke struct
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Message: "Invalid request",
            Errors: []domains.ErrorDetail{
                {Message: "Invalid Server Error", Parameter: "request body"},
            },
            Code: http.StatusInternalServerError,
        })
    }

    var errorDetails []domains.ErrorDetail

    // Ambil user berdasarkan username
    user, err := h.Usecase.GetByUsername(req.Username)
    if err != nil {
        errorDetails = append(errorDetails, domains.ErrorDetail{
            Message: "User not found",
            Parameter: "username",
        })
    } 

	if user != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
            errorDetails = append(errorDetails, domains.ErrorDetail{
                Message: "Invalid Password",
                Parameter: "password",
            })
        }
	} else {
        // Jika user tidak ditemukan, kita juga perlu memberi tahu bahwa password salah
        errorDetails = append(errorDetails, domains.ErrorDetail{
            Message: "Invalid Password",
            Parameter: "password",
        })
    }

    // Kirimkan response error jika ada
    if len(errorDetails) > 0 {
        return c.JSON(http.StatusUnauthorized, domains.Response{
            Message: "Authentication Failed",
            Data: nil,
            Errors:   errorDetails,
            Code:    http.StatusUnauthorized,
        })
    }

    // Jika tidak ada error, kembalikan response yang sukses
    return c.JSON(http.StatusOK, domains.Response{
        Message: "Valid Credentials",
        Data: map[string]interface{}{
            "username": req.Username,
        },
        Errors: nil,
        Code: http.StatusOK,
    })
}
