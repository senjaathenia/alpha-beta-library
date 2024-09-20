package delivery

import (
	"net/http"
	"project-golang-crud/domains"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AuthorHandler struct {
	AuthorUsecase domains.AuthorUsecase
}

func NewAuthorHandler(e *echo.Echo, authorUsecase domains.AuthorUsecase) {
	handler := &AuthorHandler{AuthorUsecase: authorUsecase}

	e.POST("/authors", handler.Create)
	e.GET("/authors", handler.GetAll)
	e.GET("/authors/:id", handler.GetByID)
	e.PUT("/authors/:id", handler.Update)
	e.DELETE("/authors/:id", handler.Delete)
}
func (h *AuthorHandler) Create(c echo.Context) error {
	var author domains.Author
	if err := c.Bind(&author); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if err := h.AuthorUsecase.Create(&author); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, author)
}
func (h *AuthorHandler) GetAll(c echo.Context) error {
	authors, err := h.AuthorUsecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, authors)
}
func (h *AuthorHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Id"})
	}
	author, err := h.AuthorUsecase.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if author == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Author not found"})
	}
	return c.JSON(http.StatusOK, author)
}
func (h *AuthorHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Id"})
	}

	var author domains.Author
	if err := c.Bind(&author); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	author.ID = uint(id)

	if err := h.AuthorUsecase.Update(&author); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, author)
}
func (h *AuthorHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Id"})
	}
	if err := h.AuthorUsecase.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
