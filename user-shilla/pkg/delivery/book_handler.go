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

type BookHandler struct {
	BookUsecase domains.BookUsecase
}

func NewBookHandler(e *echo.Echo, bookUsecase domains.BookUsecase) {
	handler := &BookHandler{BookUsecase: bookUsecase}

	e.POST("/books", handler.Create, middleware.JWTMiddleware("admin"))
	e.GET("/books", handler.GetAll, middleware.JWTMiddleware("admin"))
	e.GET("/books/:id", handler.GetByID, middleware.JWTMiddleware("admin"))
	e.PUT("/books/:id", handler.Update, middleware.JWTMiddleware("admin"))
	e.DELETE("/books/:id", handler.Delete, middleware.JWTMiddleware("admin"))
}

// Create godoc
// @Summary Create a new book
// @Description Create a new book with the input payload
// @Tags books
// @Accept json
// @Produce json
// @Param book body domains.Book true "Book payload"
// @Success 201 {object} domains.Response
// @Failure 400 {object} domains.Response
// @Failure 500 {object} domains.Response
// @Router /books [post]
    func (h *BookHandler) Create(c echo.Context) error {
        var book domains.Book
        if err := c.Bind(&book); err != nil {
            return c.JSON(http.StatusInternalServerError, domains.Response{
                Code:    http.StatusInternalServerError,
                Message: "Failed to parse request",
                Data:    nil,
            })
        }

        var errorDetails []domains.ErrorDetail

        // Cek apakah title kosong, jika ya tambahkan pesan error
        if book.Title == "" {
            errorDetails = append(errorDetails, domains.ErrorDetail{
                Message:  "Title is required",
                Parameter: "title",
            })
        }

		ctx := c.Request().Context()
    // Cek apakah author ada, jika tidak ada tambahkan pesan error
    if err := h.BookUsecase.AuthorExists(ctx, int(book.AuthorID)); err != nil {
        errorDetails = append(errorDetails, domains.ErrorDetail{
            Message:  "Author does not exist",
            Parameter: "author_id",
        })
    }

    // Cek apakah publisher ada, jika tidak ada tambahkan pesan error
    if err := h.BookUsecase.PublisherExists(ctx, int(book.PublisherID)); err != nil {
        errorDetails = append(errorDetails, domains.ErrorDetail{
            Message:  "Publisher does not exist",
            Parameter: "publisher_id",
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
        if err := h.BookUsecase.CheckTitleExists(ctx, book.Title, book.ID); err != nil {
            if strings.Contains(err.Error(), "exists") {
                return c.JSON(http.StatusBadRequest, domains.Response{
                    Code:    http.StatusBadRequest, 
                    Message: "Buku dengan judul ini sudah ada",
                    Data:    nil,
                })
            }
            return c.JSON(http.StatusInternalServerError, domains.Response{
                Code:    http.StatusInternalServerError,
                Message: "Internal server error",
                Data:    nil,
            })
        }

        if err := h.BookUsecase.Create(ctx, &book); err != nil {
            return c.JSON(http.StatusInternalServerError, domains.Response{
                Code:    http.StatusInternalServerError,
                Message: "Internal server error",
                Data:    nil,
            })
        }


        return c.JSON(http.StatusCreated, domains.Response{
            Code:    http.StatusCreated,
            Message: "Book created successfully",
            Data:    book,
        })
    }


// GetByID godoc
// @Summary Get book by ID
// @Description Get a single book by its ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} domains.Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [get]
func (h *BookHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()  
	books, err := h.BookUsecase.GetAll(ctx)  
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}
	return c.JSON(http.StatusOK, domains.Response{
		Code:    http.StatusOK,
		Message: "Books retrieved successfully",
		Data:    books,
	})
}

// GetByID godoc
// @Summary Get book by ID
// @Description Get a single book by its ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} domains.Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [get]
func (h *BookHandler) GetByID(c echo.Context) error {
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
    book, err := h.BookUsecase.GetByID(ctx, uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound){
            return c.JSON(http.StatusNotFound, domains.Response{
                Code: http.StatusNotFound,
                Message: "Book not found",
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
        Message: "Book retrieved successfully",
        Data: book,
    })
    }
    
// Update godoc
// @Summary Update a book
// @Description Update a book's information by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body domains.Book true "Book payload"
// @Success 200 {object} domains.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [put]
func (h *BookHandler) Update(c echo.Context) error {
    var updatedBook domains.Book
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
    existingBook, err := h.BookUsecase.GetByID(ctx, uint(id))
    if err != nil {
        return c.JSON(http.StatusNotFound, domains.Response{
            Code:    http.StatusNotFound,
            Message: "Book not found",
            Data:    nil,
        })
    }

    if err := c.Bind(&updatedBook); err != nil {
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Code:    http.StatusInternalServerError,
            Message: "Internal Server Error",
            Data:    nil,
        })
    }

    updatedBook.ID = existingBook.ID
    updatedBook.CreatedAt = existingBook.CreatedAt
    updatedBook.MaxStock = existingBook.MaxStock 

    // Periksa dan update hanya field yang diubah
    if updatedBook.Title == "" {
        updatedBook.Title = existingBook.Title
    }
    if updatedBook.AuthorID == 0 {
        updatedBook.AuthorID = existingBook.AuthorID
    }
    if updatedBook.PublisherID == 0 {
        updatedBook.PublisherID = existingBook.PublisherID
    }
    if updatedBook.Summary == "" {
        updatedBook.Summary = existingBook.Summary
    }
    if updatedBook.Stock == 0 {
        updatedBook.Stock = existingBook.Stock
    }

    // Validasi
    var errorDetails []domains.ErrorDetail 
    if updatedBook.Title == "" {
        errorDetails = append(errorDetails, domains.ErrorDetail{
            Message:  "Title is required",
            Parameter: "title",
        })
    }
    if err := h.BookUsecase.AuthorExists(ctx, int(updatedBook.AuthorID)); err != nil {
        errorDetails = append(errorDetails, domains.ErrorDetail{
            Message:  "Author does not exist",
            Parameter: "author_id",
        })
    }
    if err := h.BookUsecase.PublisherExists(ctx, int(updatedBook.PublisherID)); err != nil {
        errorDetails = append(errorDetails, domains.ErrorDetail{
            Message:  "Publisher does not exist",
            Parameter: "publisher_id",
        })
    }

    if len(errorDetails) > 0 {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Validation errors",
            Data:    nil,
            Errors:  errorDetails, 
        })
    }

    if err := h.BookUsecase.CheckTitleExists(ctx, updatedBook.Title, existingBook.ID); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Book with this title already exists",
            Data:    nil,
        })
    }

    updatedBook.UpdatedAt = time.Now()
    if err := h.BookUsecase.Update(ctx, &updatedBook); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: err.Error(),
            Data:    nil,
        })
    }

    return c.JSON(http.StatusOK, domains.Response{
        Code:    http.StatusOK,
        Message: "Book updated successfully",
        Data:    updatedBook, // Mengembalikan updatedBook yang sudah diperbarui
    })
}



// Delete godoc
// @Summary Delete a book
// @Description Delete a book by its ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [delete]
func (h *BookHandler) Delete(c echo.Context) error {
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
	existingBook, err := h.BookUsecase.GetByID(ctx, uint(id))
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

    if existingBook.DeletedAt != nil {
        return c.JSON(http.StatusNotFound, domains.Response{
            Code: http.StatusNotFound,
            Message: "Book not found",
            Data: nil,
        })
    }

    deletedAt := time.Now()
	if err := h.BookUsecase.Delete(ctx, uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
			Data: nil,
		})
	}
    existingBook.DeletedAt = &deletedAt

	return c.JSON(http.StatusOK, domains.Response{
		Code: http.StatusOK,
		Message: "Book deleted successfully",
		Data: existingBook,
	})
	}
