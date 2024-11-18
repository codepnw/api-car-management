package enghandlers

import (
	"context"
	"net/http"
	"time"

	"github.com/codepnw/go-car-management/modules/engines"
	engservices "github.com/codepnw/go-car-management/modules/engines/services"
	"github.com/gin-gonic/gin"
)

type enginHandler struct {
	service engservices.IEngineService
}

func NewEngineHandler(service engservices.IEngineService) *enginHandler {
	return &enginHandler{service: service}
}

func (h *enginHandler) GetEngineByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")

	resp, err := h.service.GetEngineByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *enginHandler) CreateEngine(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &engines.EngineRequest{}

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdEngine, err := h.service.CreateEngine(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdEngine})
}

func (h *enginHandler) UpdateEngine(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	req := &engines.EngineRequest{}

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEngine, err := h.service.UpdateEngine(ctx, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedEngine})
}

func (h *enginHandler) DeleteEngine(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")

	deletedEngine, err := h.service.DeleteEngine(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"data": deletedEngine})
}
