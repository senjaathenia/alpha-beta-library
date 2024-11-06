package delivery

import (
	"errors"
	"log"
	"net/http"
	"project-golang-crud/domains"
	"project-golang-crud/middleware"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type LoanHandler struct {
	LoanUsecase domains.LoanUsecase
}

func NewLoanHandler(e *echo.Echo, loanUsecase domains.LoanUsecase) {
	handler := &LoanHandler{LoanUsecase: loanUsecase}

	e.POST("/loans", handler.Create, middleware.JWTMiddleware("user"))
	e.GET("/loans", handler.GetAll, middleware.JWTMiddleware("admin"))
	e.GET("/loans/user", handler.GetByUsername, middleware.JWTMiddleware("admin"))
	e.PUT("/loans/:id", handler.Update, middleware.JWTMiddleware("admin"))
	e.DELETE("/loans/:id", handler.Delete, middleware.JWTMiddleware("admin"))
	e.PUT("/loans/return/:id", handler.ReturnBook, middleware.JWTMiddleware("admin"))
}

// Create godoc
// @Summary Create a new loan
// @Description Create a new loan with the input payload
// @Tags loans
// @Accept json
// @Produce json
// @Param loan body domains.BookLoans true "Loan payload"
// @Success 201 {object} domains.Response
// @Failure 400 {object} domains.Response
// @Failure 500 {object} domains.Response
// @Router /loans [post]
func (h *LoanHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var loan domains.BookLoans // Pastikan tipe data sesuai
	if err := c.Bind(&loan); err != nil {
		log.Println("Error binding loan:", err)
		return c.JSON(http.StatusBadRequest, domains.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
		})
	}

	// Validasi lainnya
	if loan.UserID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Code:    http.StatusBadRequest,
			Message: "UserID is required",
			Data:    nil,
		})
	}

	// Panggil use case untuk menyimpan data
	if err := h.LoanUsecase.Create(ctx, &loan); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create loan",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, domains.Response{
		Code:    http.StatusCreated,
		Message: "Loan created successfully",
		Data:    loan,
	})
}

// GetAll godoc
// @Summary Get all loans
// @Description Retrieve all loans
// @Tags loans
// @Produce json
// @Success 200 {object} domains.Response
// @Failure 500 {object} domains.Response
// @Router /loans [get]
func (h *LoanHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	loans, err := h.LoanUsecase.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}
	return c.JSON(http.StatusOK, domains.Response{
		Code:    http.StatusOK,
		Message: "BookLoans retrieved successfully",
		Data:    loans,
	})
}

// GetByID godoc
// @Summary Get loan by ID
// @Description Get a single loan by its ID
// @Tags loans
// @Produce json
// @Param id path int true "Loan ID"
// @Success 200 {object} domains.BookLoans
// @Failure 400 {object} domains.Response
// @Failure 404 {object} domains.Response
// @Failure 500 {object} domains.Response
// @Router /loans/{id} [get]
func (h *LoanHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid Id",
			Data:    nil,
		})
	}
	ctx := c.Request().Context()
	loan, err := h.LoanUsecase.GetByID(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, domains.Response{
				Code:    http.StatusNotFound,
				Message: "Loan not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data:    nil,
		})
	}
	return c.JSON(http.StatusOK, domains.Response{
		Code:    http.StatusOK,
		Message: "Loan retrieved successfully",
		Data:    loan,
	})
}

// Update godoc
// @Summary Update a loan
// @Description Update a loan's information by its ID
// @Tags loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param loan body domains.BookLoans true "Loan payload"
// @Success 200 {object} domains.BookLoans
// @Failure 400 {object} domains.Response
// @Failure 500 {object} domains.Response
// @Router /loans/{id} [put]
func (h *LoanHandler) Update(c echo.Context) error {
	var loan domains.BookLoans
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
	existingLoan, err := h.LoanUsecase.GetByID(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, domains.Response{
			Code:    http.StatusNotFound,
			Message: "Loan not found",
			Data:    nil,
		})
	}

	// Bind data from request body to loan
	if err := c.Bind(&loan); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data:    nil,
		})
	}
	loan.ID = existingLoan.ID // Set the ID for the updated loan

	var errorDetails []domains.ErrorDetail

	// Validasi BookID
	if loan.BookID == 0 {
		errorDetails = append(errorDetails, domains.ErrorDetail{
			Message:   "BookID is required",
			Parameter: "book_id",
		})
	}

	// Validasi UserID
	if loan.UserID == uuid.Nil {
		errorDetails = append(errorDetails, domains.ErrorDetail{
			Message:   "UserID is required",
			Parameter: "user_id",
		})
	}

	// Jika ada error, kembalikan semua error sebagai response
	if len(errorDetails) > 0 {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation errors",
			Data:    nil,
			Errors:  errorDetails,
		})
	}

	if err := h.LoanUsecase.Update(ctx, &loan); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, domains.Response{
		Code:    http.StatusOK,
		Message: "Loan updated successfully",
		Data:    loan,
	})
}

// Delete godoc
// @Summary Delete a loan
// @Description Delete a loan by its ID
// @Tags loans
// @Produce json
// @Param id path int true "Loan ID"
// @Success 204
// @Failure 400 {object} domains.Response
// @Failure 500 {object} domains.Response
// @Router /loans/{id} [delete]
func (h *LoanHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid Id",
			Data:    nil,
		})
	}
	ctx := c.Request().Context()
	existingLoan, err := h.LoanUsecase.GetByID(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, domains.Response{
				Code:    http.StatusNotFound,
				Message: "Loan not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data:    nil,
		})
	}

	if err := h.LoanUsecase.Delete(ctx, uint(existingLoan.ID)); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data:    nil,
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ReturnBook godoc
// @Summary Return a borrowed book
// @Description Process the return of a borrowed book and calculate fines if any
// @Tags loans
// @Param id path int true "Loan ID"
// @Success 200 {object} domains.Response
// @Failure 400 {object} domains.Response
// @Failure 500 {object} domains.Response
// @Router /loans/return/{id} [put]
// ReturnBook godoc
// @Summary Return a borrowed book
// @Description Process the return of a borrowed book and calculate fines if any
// @Tags loans
// @Param id path int true "Loan ID"
// @Success 200 {object} domains.Response
// @Failure 400 {object} domains.Response
// @Failure 500 {object} domains.Response
// @Router /loans/return/{id} [put]
func (h *LoanHandler) ReturnBook(c echo.Context) error {
	const lateFeeRate = 5000 // Tarif denda per hari, misalnya 5000 IDR per hari
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid Loan ID",
            Data:    nil,
        })
    }

    // Cek payload JSON
    var request struct {
        ReturnDate string `json:"return_date"`
    }

    if err := c.Bind(&request); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid Request Body",
            Data:    nil,
        })
    }

    // Konversi string ke time.Time
    returnDate, err := time.Parse(time.RFC3339, request.ReturnDate)
    if err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid Return Date: " + err.Error(),
            Data:    nil,
        })
    }

    ctx := c.Request().Context()
    
    // Ambil detail peminjaman untuk menghitung denda
    bookLoan, err := h.LoanUsecase.GetByID(ctx, uint(id))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Code:    http.StatusInternalServerError,
            Message: err.Error(),
            Data:    nil,
        })
    }

    // Hitung late fee jika perlu
    lateFee := 0
    if returnDate.After(bookLoan.DueDate) {
        lateDays := int(returnDate.Sub(bookLoan.DueDate).Hours() / 24)
        lateFee = lateDays * lateFeeRate // misalkan lateFeeRate sudah didefinisikan
    }

    err = h.LoanUsecase.Return(ctx, uint(id), returnDate)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Code:    http.StatusInternalServerError,
            Message: err.Error(),
            Data:    nil,
        })
    }

    return c.JSON(http.StatusOK, domains.Response{
        Code:    http.StatusOK,
        Message: "Book returned successfully",
        Data: map[string]interface{}{
            "late_fee": lateFee,
        },
    })
}

func (h *LoanHandler) GetByUsername(c echo.Context) error {
    var request struct {
        Username string `json:"username"`
    }

    if err := c.Bind(&request); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Invalid request format",
            Data:    nil,
            Errors: []domains.ErrorDetail{
                {Message: "Invalid request format", Parameter: "username"},
            },
        })
    }

    if request.Username == "" {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Username is required",
            Data:    nil,
            Errors: []domains.ErrorDetail{
                {Message: "Username is required", Parameter: "username"},
            },
        })
    }

    ctx := c.Request().Context()
    loans, err := h.LoanUsecase.GetByUsernameLoans(ctx, request.Username)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.JSON(http.StatusNotFound, domains.Response{
                Code:    http.StatusNotFound,
                Message: "No loans found for the specified username",
                Data:    nil,
            })
        }
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Code:    http.StatusInternalServerError,
            Message: "Internal Server Error",
            Data:    nil,
        })
    }

    return c.JSON(http.StatusOK, domains.Response{
        Code:    http.StatusOK,
        Message: "Loans retrieved successfully",
        Data:    loans,
    })
}

