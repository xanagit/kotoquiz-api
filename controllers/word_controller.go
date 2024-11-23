package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type WordController interface {
	ReadWord(c *gin.Context)
	CreateWord(c *gin.Context)
	UpdateWord(c *gin.Context)
	DeleteWord(c *gin.Context)
}

type WordControllerImpl struct {
	Service services.WordService
}

func (s *WordControllerImpl) ReadWord(c *gin.Context) {
	id := c.Param("id")
	word, err := s.Service.ReadWord(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, word)
}

func (s *WordControllerImpl) CreateWord(c *gin.Context) {
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	word.ID = uuid.Nil
	if err := s.Service.CreateWord(&word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, word)
}

func (s *WordControllerImpl) UpdateWord(c *gin.Context) {
	id := c.Param("id")
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	word.ID = FromStrToUuid(id)
	if err := s.Service.UpdateWord(&word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, word)
}

func (s *WordControllerImpl) DeleteWord(c *gin.Context) {
	id := c.Param("id")
	if err := s.Service.DeleteWord(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
