package delivery

import (
	"errors"
	"net/http"
	"project-golang-crud/domains"
	"project-golang-crud/middleware"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type BookRequestHandler struct {
	usecase domains.BookRequestUsecase
}

func NewBookRequestHandler(e *echo.Echo, uc domains.BookRequestUsecase) {
	handler := &BookRequestHandler{usecase: uc}
	e.POST("/book-requests", handler.Create, middleware.JWTMiddleware("user"))
	e.GET("/book-requests", handler.GetAll, middleware.JWTMiddleware("admin"))
	e.GET("/book-requests/:id", handler.GetByID, middleware.JWTMiddleware("admin"))
	e.PUT("/book-requests/:id", handler.Update, middleware.JWTMiddleware("admin"))
	e.DELETE("/book-requests/:id", handler.Delete, middleware.JWTMiddleware("admin"))
	e.POST("/book-requests/:id/approve", handler.ApproveOrReject, middleware.JWTMiddleware("admin"))
	e.GET("/book-requests/user", handler.GetByUsername, middleware.JWTMiddleware("admin"))
}
func (h *BookRequestHandler) GetByUsername(c echo.Context) error {
    // Memastikan metode yang digunakan adalah GET
    if c.Request().Method != http.MethodGet {
        return c.JSON(http.StatusMethodNotAllowed, domains.Response{
            Message: "Method Not Allowed",
            Data:    nil,
            Errors: []domains.ErrorDetail{
                {Message: "This endpoint does not support the requested method", Parameter: "method"},
            },
            Code: http.StatusMethodNotAllowed,
        })
    }

    // Baca body request
    var requestBody struct {
        Username string `json:"username"`
    }

    // Mengambil body JSON
    if err := c.Bind(&requestBody); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Message: "Invalid request body",
            Data:    nil,
            Errors: []domains.ErrorDetail{
                {Message: "Invalid JSON", Parameter: "body"},
            },
            Code: http.StatusBadRequest,
        })
    }

    // Pastikan username tidak kosong
    if requestBody.Username == "" {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Message: "Username cannot be empty",
            Data:    nil,
            Errors: []domains.ErrorDetail{
                {Message: "Username cannot be empty", Parameter: "username"},
            },
            Code: http.StatusBadRequest,
        })
    }

    // Ambil data user dan book_requests berdasarkan username
    userWithLoans, err := h.usecase.GetByUsernameRequest(c.Request().Context(), requestBody.Username)
    if err != nil {
        if errors.Is(err, domains.ErrUserNotFound) {
            return c.JSON(http.StatusNotFound, domains.Response{
                Message: "User not found",
                Data:    nil,
                Errors: []domains.ErrorDetail{
                    {Message: "User with the given username does not exist", Parameter: "username"},
                },
                Code: http.StatusNotFound,
            })
        }
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Message: "Internal Server Error",
            Data:    nil,
            Errors:  nil,
            Code:    http.StatusInternalServerError,
        })
    }

    // Kembalikan respons dengan data book_requests dan username
    return c.JSON(http.StatusOK, domains.Response{
        Message: "User Request Found",
        Data:    userWithLoans.BookRequest,
        Errors:  nil,
        Code:    http.StatusOK,
    })
}


func (h *BookRequestHandler) Create(c echo.Context) error {
	var request domains.BookRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Message: "Invalid request format",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: err.Error(), Parameter: "body"},
			},
			Code: http.StatusBadRequest,
		})
	}

	// Cek ketersediaan buku
	book, err := h.usecase.GetBookByID(c.Request().Context(), request.BookID) // Ambil informasi buku
	if err != nil {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Message: "Book is not available",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: err.Error(), Parameter: "book_id"},
			},
			Code: http.StatusBadRequest,
		})
	}

	// Cek stok buku
	if book.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Message: "Stock not available",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: "The book stock is not available", Parameter: "book_id"},
			},
			Code: http.StatusBadRequest,
		})
	}

	// Jika stok tersedia, lakukan pengurangan stok
	if err := h.usecase.DecreaseStock(c.Request().Context(), request.BookID); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Message: "Failed to decrease stock",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: err.Error(), Parameter: "stock"},
			},
			Code: http.StatusInternalServerError,
		})
	}

	// Buat permintaan peminjaman
	if err := h.usecase.Create(c.Request().Context(), &request); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Message: "Failed to create book request",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: err.Error(), Parameter: "request"},
			},
			Code: http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusCreated, domains.Response{
		Message: "Book request created successfully",
		Data:    request,
		Errors:  nil,
		Code:    http.StatusCreated,
	})
}

func (h *BookRequestHandler) GetAll(c echo.Context) error {
	requests, err := h.usecase.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Message: "Failed to retrieve book requests",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: err.Error(), Parameter: "requests"},
			},
			Code: http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, domains.Response{
		Message: "Book requests retrieved successfully",
		Data:    requests,
		Errors:  nil,
		Code:    http.StatusOK,
	})
}

func (h *BookRequestHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Message: "Invalid ID format",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: "ID must be a number", Parameter: "id"},
			},
			Code: http.StatusBadRequest,
		})
	}

	request, err := h.usecase.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, domains.Response{
			Message: "Book request not found",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: "No request found with the given ID", Parameter: "id"},
			},
			Code: http.StatusNotFound,
		})
	}

	return c.JSON(http.StatusOK, domains.Response{
		Message: "Book request retrieved successfully",
		Data:    request,
		Errors:  nil,
		Code:    http.StatusOK,
	})
}

func (h *BookRequestHandler) Update(c echo.Context) error {
	var request domains.BookRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Message: "Invalid request body",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: err.Error(), Parameter: "body"},
			},
			Code: http.StatusBadRequest,
		})
	}

	id, _ := strconv.Atoi(c.Param("id"))
	request.ID = uint(id)

	if err := h.usecase.Update(c.Request().Context(), &request); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Message: "Failed to update book request",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: err.Error(), Parameter: "request"},
			},
			Code: http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, domains.Response{
		Message: "Book request updated successfully",
		Data:    request,
		Errors:  nil,
		Code:    http.StatusOK,
	})
}

func (h *BookRequestHandler) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.usecase.Delete(c.Request().Context(), uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Message: "Failed to delete book request",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: err.Error(), Parameter: "request"},
			},
			Code: http.StatusInternalServerError,
		})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *BookRequestHandler) ApproveOrReject(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Message: "Invalid ID format",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: "ID must be a number", Parameter: "id"},
			},
			Code: http.StatusBadRequest,
		})
	}

	var req domains.LoanUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Message: "Invalid request body",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: "Invalid request body: " + err.Error(), Parameter: "body"},
			},
			Code: http.StatusBadRequest,
		})
	}

	c.Logger().Infof("Parsed request: %+v", req)

	if !req.Approved && req.RejectionReason == "" {
		return c.JSON(http.StatusBadRequest, domains.Response{
			Message: "Rejection reason is required",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: "Rejection reason is required if the request is not approved", Parameter: "rejection_reason"},
			},
			Code: http.StatusBadRequest,
		})
	}

	dueDate, err := h.usecase.ApproveOrReject(c.Request().Context(), uint(id), req.Approved, req.RejectionReason)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Message: "Failed to update request status",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: err.Error(), Parameter: "request"},
			},
			Code: http.StatusInternalServerError,
		})
	}

	updatedRequest, err := h.usecase.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Message: "Failed to retrieve updated request",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: "Failed to retrieve updated request", Parameter: "request"},
			},
			Code: http.StatusInternalServerError,
		})
	}

	loc, err := time.LoadLocation("Asia/Jakarta") // Set to the desired time zone
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Message: "Failed to load location",
			Data:    nil,
			Errors: []domains.ErrorDetail{
				{Message: err.Error(), Parameter: "location"},
			},
			Code: http.StatusInternalServerError,
		})
	}

	// Formatting dates in the response to include the correct time zone
	updatedRequest.RequestDate = updatedRequest.RequestDate.In(loc)
	updatedRequest.CreatedAt = updatedRequest.CreatedAt.In(loc)
	updatedRequest.UpdatedAt = updatedRequest.UpdatedAt.In(loc)

	response := domains.Response{
		Message: "Request status updated successfully",
		Data:    updatedRequest, // Menggunakan updatedRequest yang telah diisi
		Errors:  nil,
		Code:    http.StatusOK,
	}
	
	// Pastikan username diambil dari updatedRequest
	response.Data = map[string]interface{}{
		"request": updatedRequest,
		"due_date": dueDate.In(loc).Format("2006-01-02T15:04:05.999999999+07:00"),
		"username": updatedRequest.Username, // Pastikan username diambil dari updatedRequest
	}	

	return c.JSON(http.StatusOK, response)
}
