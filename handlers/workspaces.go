package handlers

import (
	"net/http"

	"github.com/freekobie/hazel/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateWorkspace(c *gin.Context) {
	var input struct {
		Name        string    `json:"name" binding:"required"`
		Description string    `json:"description"`
		UserID      uuid.UUID `json:"userId" binding:"required,uuid"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ws := &models.Workspace{
		Name:        input.Name,
		Description: input.Description,
		User:        &models.User{Id: input.UserID},
	}

	err = h.wss.NewWorkspace(c.Request.Context(), ws)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ws)
}
