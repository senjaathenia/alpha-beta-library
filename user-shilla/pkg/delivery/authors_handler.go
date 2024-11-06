package delivery

import (
	"errors"
	"net/http"
	"project-golang-crud/domains"
	"project-golang-crud/middleware"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthorHandler struct {
	AuthorUsecase domains.AuthorUsecase
}

func NewAuthorHandler(e *echo.Echo, authorUsecase domains.AuthorUsecase) {
	handler := &AuthorHandler{AuthorUsecase: authorUsecase}

	e.POST("/authors", handler.Create, middleware.JWTMiddleware("admin"))
	e.GET("/authors", handler.GetAll, middleware.JWTMiddleware("admin"))
	e.GET("/authors/:id", handler.GetByID, middleware.JWTMiddleware("admin"))
	e.PUT("/authors/:id", handler.Update, middleware.JWTMiddleware("admin"))
	e.DELETE("/authors/:id", handler.Delete, middleware.JWTMiddleware("admin"))
}

func (h *AuthorHandler) Create(c echo.Context) error {
	var author domains.Author
	if err := c.Bind(&author); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse request",
			Data:    nil,
		})
	}

	var errorDetails []domains.ErrorDetail

	// Cek apakah title kosong, jika ya tambahkan pesan error
	if author.Name == "" {
		errorDetails = append(errorDetails, domains.ErrorDetail{
			Message:  "Name is required",
			Parameter: "name",
		})
	}

	// Jika ada error, kembalikan semua error sebagai response
	if len(errorDetails) > 0 {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation errors",
			Data:    nil,
			Errors: errorDetails,
		})
	}

	// Memanggil usecase untuk menghandle logika bisnis (buku hanya dibuat setelah validasi author dan publisher)
	ctx := c.Request().Context()
	if err := h.AuthorUsecase.CheckNameExists(ctx, author.Name, author.ID); err != nil {
		if strings.Contains(err.Error(), "exists") {
			return c.JSON(http.StatusBadRequest, domains.Response{
				Code:    http.StatusBadRequest, 
				Message: "Author with this name is already exist",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	if err := h.AuthorUsecase.Create(ctx, &author); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Data:    nil,
		})
	}


	return c.JSON(http.StatusCreated, domains.Response{
		Code:    http.StatusCreated,
		Message: "Book created successfully",
		Data:    author,
	})
}


func (h *AuthorHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()  
	authors, err := h.AuthorUsecase.GetAll(ctx)  
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}
	return c.JSON(http.StatusOK, domains.Response{
		Code:    http.StatusOK,
		Message: "Author retrieved successfully",
		Data:    authors,
	})
}

func (h *AuthorHandler) GetByID(c echo.Context) error {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
             Code: http.StatusBadRequest,
             Message: "Invalid Id",
             Data: nil,
        })
    }
    ctx := c.Request().Context()
    book, err := h.AuthorUsecase.GetByID(ctx, uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound){
            return c.JSON(http.StatusNotFound, domains.Response{
                Code: http.StatusNotFound,
                Message: "Author not found",
                Data: nil,
            })
        }
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Code: http.StatusInternalServerError,
            Message: "Internal Server Error",
            Data: nil,
        })
    }
    return c.JSON(http.StatusOK, domains.Response{
        Code: http.StatusOK,
        Message: "Author retrieved successfully",
        Data: book,
    })
    }
 
func (h *AuthorHandler) Update(c echo.Context) error {
    var author *domains.Author
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid Id",
            Data:    nil,
        })
    }   

    ctx := c.Request().Context()
    author, err = h.AuthorUsecase.GetByID(ctx, uint(id))
    if err != nil {
        return c.JSON(http.StatusNotFound, domains.Response{ // Mengubah status menjadi NotFound
            Code:    http.StatusNotFound,
            Message: "Author not found",
            Data:    nil,
        })
    }

    // Bind data from request body to updatedBook
    var updatedAuthor domains.Author
    if err := c.Bind(&updatedAuthor); err != nil {
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Code:    http.StatusInternalServerError,
            Message: "Internal Server Error",
            Data:    nil,
        })
    }
    updatedAuthor.ID = uint(id) // Set the ID for the updated book
    updatedAuthor.CreatedAt = author.CreatedAt

    var errorDetails []domains.ErrorDetail // Menggunakan ErrorDetail

    // Cek apakah title kosong
    if updatedAuthor.Name == "" {
        errorDetails = append(errorDetails, domains.ErrorDetail{
            Message:  "Name is required",
            Parameter: "name",
        })
    }
 
    // Jika ada error, kembalikan semua error sebagai response
    if len(errorDetails) > 0 {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Validation errors",
            Data:    nil,
            Errors:  errorDetails, // Tambahkan detail error ke response
        })
    }

    if err := h.AuthorUsecase.CheckNameExists(ctx, updatedAuthor.Name, uint(id)); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Author with this name already exists",
            Data:    nil,
        })
    }

    updatedAuthor.UpdatedAt = time.Now()
    // Call the use case to update the book
    if err := h.AuthorUsecase.Update(ctx, &updatedAuthor); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: err.Error(),
            Data:    nil,
        })
    }

    return c.JSON(http.StatusOK, domains.Response{
        Code:    http.StatusOK,
        Message: "Author updated successfully",
        Data:    updatedAuthor,
    })
}

func (h *AuthorHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Code: http.StatusBadRequest,
			Message: "Invalid Id",
			Data: nil,
		})
	}
    ctx := c.Request().Context()
	existingAuthor, err := h.AuthorUsecase.GetByID(ctx, uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.JSON(http.StatusNotFound, domains.Response{
                Code:    http.StatusNotFound,
                Message: "Book not found",
                Data:    nil,
            })
        }
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Code:    http.StatusInternalServerError,
            Message: "Internal server error",
            Data:    nil,
        })
    }

    if existingAuthor.DeletedAt != nil {
        return c.JSON(http.StatusNotFound, domains.Response{
            Code: http.StatusNotFound,
            Message: "Author not found",
            Data: nil,
        })
    }

    deletedAt := time.Now()
	if err := h.AuthorUsecase.Delete(ctx, uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
			Data: nil,
		})
	}
    existingAuthor.DeletedAt = &deletedAt

	return c.JSON(http.StatusOK, domains.Response{
		Code: http.StatusOK,
		Message: "Book deleted successfully",
		Data: existingAuthor,
	})
	}