// Package controllers implements HTTP handlers for the application API endpoints
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

// WordDtoController defines the interface for word DTO-related HTTP endpoints
// It provides methods to retrieve words in a format optimized for client applications
type WordDtoController interface {
	// ListDtoWords handles GET requests to retrieve words by IDs in DTO format
	ListWordsIDs(c *gin.Context)
	// ReadDtoWord handles GET requests to retrieve a specific word by ID in DTO format
	ReadDtoWord(c *gin.Context)
	// ListWordsIDs handles GET requests to retrieve a list of word IDs with filtering
	ListDtoWords(c *gin.Context)
}

// WordDtoControllerImpl implements the WordDtoController interface
// It depends on the WordDtoService for business logic operations
type WordDtoControllerImpl struct {
	WordDtoService services.WordDtoService
}

// Make sure that WordDtoControllerImpl implements WordDtoController
var _ WordDtoController = (*WordDtoControllerImpl)(nil)

// ListWordsIDs handles GET requests to retrieve a list of word IDs with filtering
// It supports filtering by tags and levels, with pagination
//
// Query Parameters:
//   - tags: Comma-separated list of tag IDs to filter by
//   - levels: Comma-separated list of level IDs to filter by
//   - limit: Maximum number of results to return (default: 30)
//   - offset: Number of results to skip (default: 0)
//
// Responses:
//   - 200 OK with an array of word IDs on success
//   - 400 Bad Request if the parameters are invalid
//   - 500 Internal Server Error if a server error occurs
func (s *WordDtoControllerImpl) ListWordsIDs(c *gin.Context) {
	tagIds := getQueryParamList(c, "tags")
	levelNameIds := getQueryParamList(c, "levelNames")
	nb, _ := getQueryParamInt(c, "nb", DefaultQpVals.NbIdsList)

	userIDStr := c.Query("userId")
	userID, ok := parseUUID(userIDStr)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid USer id format"})
		return
	}

	wordIdsList, err := s.WordDtoService.ListWordsIDs(userID, tagIds, levelNameIds, nb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wordIdsList)
}

// ListDtoWords handles GET requests to retrieve words by IDs in DTO format
// It accepts query parameters for IDs and language
//
// Query Parameters:
//   - ids: Comma-separated list of word IDs to retrieve
//   - lang: Language code for translations (default: "en")
//
// Responses:
//   - 200 OK with an array of word DTOs on success
//   - 400 Bad Request if the parameters are invalid
//   - 500 Internal Server Error if a server error occurs
func (s *WordDtoControllerImpl) ListDtoWords(c *gin.Context) {
	rawIds := getQueryParamList(c, "ids") // Récupère les IDs depuis le paramètre de requête
	ids, ok := parseUUIDs(rawIds)
	if !ok {
		return
	}
	lang := getQueryParamLang(c)

	var words []*dto.WordDTO
	var err error

	if len(ids) > 0 {
		words, err = s.WordDtoService.ListWordsDtoByIDs(ids, lang)
	} else {
		words = []*dto.WordDTO{}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, words)
}

// ReadDtoWord handles GET requests to retrieve a specific word by ID in DTO format
// It accepts a URL parameter for the word ID and a query parameter for language
//
// URL Parameters:
//   - id: Word ID to retrieve
//
// Query Parameters:
//   - lang: Language code for translations (default: "en")
//
// Responses:
//   - 200 OK with the word DTO on success
//   - 400 Bad Request if the ID is invalid
//   - 404 Not Found if no word with the given ID exists
//   - 500 Internal Server Error if a server error occurs
func (s *WordDtoControllerImpl) ReadDtoWord(c *gin.Context) {
	rawId := c.Param("id")
	id, ok := parseUUID(rawId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}
	lang := getQueryParamLang(c)

	wordDto, err := s.WordDtoService.ReadWord(id, lang)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wordDto)
}
