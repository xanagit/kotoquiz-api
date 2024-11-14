package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xanagit/kotoquiz-api/models"
	"github.com/xanagit/kotoquiz-api/services"
	"log"
	"net/http"
)

type WordController interface {
	GetWords() ([]*models.Word, error)
	GetWordByID(id string) (*models.Word, error)
	CreateWord(word *models.Word) error
	UpdateWord(word *models.Word) error
	DeleteWord(id string) error
}

type WordControllerImpl struct {
	Service services.WordService
}

func (s *WordControllerImpl) GetWords(c *gin.Context) {
	words, err := s.Service.GetWords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, words)
}

func (s *WordControllerImpl) GetWordByID(c *gin.Context) {
	id := c.Param("id")
	word, err := s.Service.GetWordByID(id)
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
	id := c.Param("id")
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	word.ID = fromStrToUuid(id)
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

func fromStrToUuid(id string) uuid.UUID {
	parsed, err := uuid.Parse(id)
	if err != nil {
		log.Fatalf("Invalid UUID format: %v", err)
	}
	return parsed
}
