package handler

import (
	"effective/internal/model"
	"effective/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PeopleHandler struct {
	service service.PeopleService
}

func NewPeopleHandler(service service.PeopleService) *PeopleHandler {
	return &PeopleHandler{service: service}
}

func (h *PeopleHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/people")
	{
		group.POST("", h.AddPerson)
		group.GET("", h.GetPeople)
		group.PUT("/:id", h.UpdatePerson)
		group.DELETE("/:id", h.DeletePerson)
	}
}

// AddPerson godoc
// @Summary Add a new person
// @Description Add a new person with enrichment
// @Tags people
// @Accept json
// @Produce json
// @Param input body model.InputPerson true "Person data"
// @Success 201 {object} model.Person
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/people [post]
func (h *PeopleHandler) AddPerson(c *gin.Context) {
	var input model.InputPerson
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person, err := h.service.AddPerson(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, person)
}

// GetPeople godoc
// @Summary Get people
// @Description Get people with filtering and pagination
// @Tags people
// @Produce json
// @Param name query string false "Filter by name"
// @Param surname query string false "Filter by surname"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} model.Person
// @Failure 500 {object} map[string]string
// @Router /api/people [get]
func (h *PeopleHandler) GetPeople(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filters := make(map[string]string)
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}

	if surname := c.Query("surname"); surname != "" {
		filters["surname"] = surname
	}

	people, err := h.service.GetPeople(filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, people)
}

// UpdatePerson godoc
// @Summary Update a person
// @Description Update person data
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param input body map[string]interface{} true "Update data"
// @Success 200 {object} model.Person
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/people/{id} [put]
func (h *PeopleHandler) UpdatePerson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person, err := h.service.UpdatePerson(uint(id), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

// DeletePerson godoc
// @Summary Delete a person
// @Description Delete a person by ID
// @Tags people
// @Param id path int true "Person ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/people/{id} [delete]
func (h *PeopleHandler) DeletePerson(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.service.DeletePerson(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
