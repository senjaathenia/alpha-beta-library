package delivery

import (
	"net/http"
	"project-golang-crud/domains"
	"strconv"
	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	BookUsecase domains.BookUsecase
}

func NewBookHandler(e *echo.Echo, bookUsecase domains.BookUsecase) {
	handler := &BookHandler{BookUsecase: bookUsecase}

	e.POST("/books", handler.Create)
	e.GET("/books", handler.GetAll)
	e.GET("/books/:id", handler.GetByID)
	e.PUT("/books/:id", handler.Update)
	e.DELETE("/books/:id", handler.Delete)
}
func (h *BookHandler) Create(c echo.Context) error {
	var book domains.Book
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.BookUsecase.Create(&book); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, book) // Status 201 untuk create
}

func (h *BookHandler) GetAll(c echo.Context) error {
	books, err := h.BookUsecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	book, err := h.BookUsecase.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if book == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Book not found"})
	}

	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var book domains.Book
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	book.ID = uint(id)

	if err := h.BookUsecase.Update(&book); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	if err := h.BookUsecase.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent) // Status 204 untuk delete
}
