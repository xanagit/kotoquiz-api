package controllers

import (
	"github.com/gin-gonic/gin"
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
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
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

	if err := s.Service.CreateWord(&word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, word)
}

func (s *WordControllerImpl) UpdateWord(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	word.ID = id
	if err := s.Service.UpdateWord(&word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, word)
}

func (s *WordControllerImpl) DeleteWord(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	if err := s.Service.DeleteWord(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
