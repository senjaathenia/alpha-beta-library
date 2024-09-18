package delivery

import (
	"net/http"
	"project-golang-crud/domains"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PublisherHandler struct {
	PublisherUsecase domains.PublisherUsecase
}

func NewPublisherHandler(e *echo.Echo, publisherUsecase domains.PublisherUsecase) {
	handler := &PublisherHandler{PublisherUsecase: publisherUsecase}

	e.POST("/publishers", handler.Create)
	e.GET("/publishers", handler.GetAll)
	e.GET("/publishers/:id", handler.GetByID)
	e.PUT("/publishers/:id", handler.Update)
	e.DELETE("/publishers/:id", handler.Delete)
}
func (h *PublisherHandler) Create(c echo.Context) error {
	var publisher domains.Publisher
	if err := c.Bind(&publisher); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Input"})
	}
	if err := h.PublisherUsecase.Create(&publisher); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, publisher)
}
func (h *PublisherHandler) GetAll(c echo.Context) error {
	publishers, err := h.PublisherUsecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, publishers)
}
func (h *PublisherHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Id"})
	}
	publisher, err := h.PublisherUsecase.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if publisher == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Publisher not found"})
	}
	return c.JSON(http.StatusOK, publisher)
}
func (h *PublisherHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Id"})
	}

	var publisher domains.Publisher
	if err := c.Bind(&publisher); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	publisher.ID = uint(id)

	if err := h.PublisherUsecase.Update(&publisher); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, publisher)
}
func (h *PublisherHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Id"})
	}
	if err := h.PublisherUsecase.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
