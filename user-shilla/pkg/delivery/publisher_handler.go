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

type PublisherHandler struct {
	PublisherUsecase domains.PublisherUsecase
}

func NewPublisherHandler(e *echo.Echo, publisherUsecase domains.PublisherUsecase) {
	handler := &PublisherHandler{PublisherUsecase: publisherUsecase}

	e.POST("/publishers", handler.Create, middleware.JWTMiddleware("admin"))
	e.GET("/publishers", handler.GetAll, middleware.JWTMiddleware("admin"))
	e.GET("/publishers/:id", handler.GetByID, middleware.JWTMiddleware("admin"))
	e.PUT("/publishers/:id", handler.Update, middleware.JWTMiddleware("admin"))
	e.DELETE("/publishers/:id", handler.Delete, middleware.JWTMiddleware("admin"))
}

func (h *PublisherHandler) Create(c echo.Context) error {
	var publisher domains.Publisher
	if err := c.Bind(&publisher); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse request",
			Data:    nil,
		})
	}

	var errorDetails []domains.ErrorDetail

	// Cek apakah title kosong, jika ya tambahkan pesan error
	if publisher.Name == "" {
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
	if err := h.PublisherUsecase.CheckNameExists(ctx, publisher.Name, publisher.ID); err != nil {
		if strings.Contains(err.Error(), "exists") {
			return c.JSON(http.StatusBadRequest, domains.Response{
				Code:    http.StatusBadRequest, 
				Message: "Publisher with this name is already exist",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Data:    nil,
		})
	}

	if err := h.PublisherUsecase.Create(ctx, &publisher); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Data:    nil,
		})
	}


	return c.JSON(http.StatusCreated, domains.Response{
		Code:    http.StatusCreated,
		Message: "Publisher created successfully",
		Data:    publisher,
	})
}


func (h *PublisherHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()  
	publisher, err := h.PublisherUsecase.GetAll(ctx)  
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}
	return c.JSON(http.StatusOK, domains.Response{
		Code:    http.StatusOK,
		Message: "Publisher retrieved successfully",
		Data:    publisher,
	})
}

func (h *PublisherHandler) GetByID(c echo.Context) error {
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
    publisher, err := h.PublisherUsecase.GetByID(ctx, uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound){
            return c.JSON(http.StatusNotFound, domains.Response{
                Code: http.StatusNotFound,
                Message: "Publisher not found",
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
        Message: "Publisher retrieved successfully",
        Data: publisher,
    })
    }

 
func (h *PublisherHandler) Update(c echo.Context) error {
    var publisher *domains.Publisher
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
    publisher, err = h.PublisherUsecase.GetByID(ctx, uint(id))
    if err != nil {
        return c.JSON(http.StatusNotFound, domains.Response{ // Mengubah status menjadi NotFound
            Code:    http.StatusNotFound,
            Message: "Publisher not found",
            Data:    nil,
        })
    }

    // Bind data from request body to updatedBook
    var updatedPublisher domains.Publisher
    if err := c.Bind(&updatedPublisher); err != nil {
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Code:    http.StatusInternalServerError,
            Message: "Internal Server Error",
            Data:    nil,
        })
    }
    updatedPublisher.ID = uint(id) // Set the ID for the updated book
    updatedPublisher.CreatedAt = publisher.CreatedAt

    var errorDetails []domains.ErrorDetail // Menggunakan ErrorDetail

    // Cek apakah title kosong
    if updatedPublisher.Name == "" {
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

    if err := h.PublisherUsecase.CheckNameExists(ctx, updatedPublisher.Name, uint(id)); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: "Publisher with this name already exists",
            Data:    nil,
        })
    }

    updatedPublisher.UpdatedAt = time.Now()
    // Call the use case to update the book
    if err := h.PublisherUsecase.Update(ctx, &updatedPublisher); err != nil {
        return c.JSON(http.StatusBadRequest, domains.Response{
            Code:    http.StatusBadRequest,
            Message: err.Error(),
            Data:    nil,
        })
    }

    return c.JSON(http.StatusOK, domains.Response{
        Code:    http.StatusOK,
        Message: "Author updated successfully",
        Data:    updatedPublisher,
    })
}

func (h *PublisherHandler) Delete(c echo.Context) error {
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
	existingPublisher, err := h.PublisherUsecase.GetByID(ctx, uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.JSON(http.StatusNotFound, domains.Response{
                Code:    http.StatusNotFound,
                Message: "Publisher not found",
                Data:    nil,
            })
        }
        return c.JSON(http.StatusInternalServerError, domains.Response{
            Code:    http.StatusInternalServerError,
            Message: "Internal server error",
            Data:    nil,
        })
    }

    if existingPublisher.DeletedAt != nil {
        return c.JSON(http.StatusNotFound, domains.Response{
            Code: http.StatusNotFound,
            Message: "Publisher not found",
            Data: nil,
        })
    }

    deletedAt := time.Now()
	if err := h.PublisherUsecase.Delete(ctx, uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, domains.Response{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
			Data: nil,
		})
	}
    existingPublisher.DeletedAt = &deletedAt

	return c.JSON(http.StatusOK, domains.Response{
		Code: http.StatusOK,
		Message: "Publisher deleted successfully",
		Data: existingPublisher,
	})
	}