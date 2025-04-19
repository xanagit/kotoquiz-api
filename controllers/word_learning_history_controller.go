// Package controllers provides HTTP handlers for the application's API endpoints
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/middlewares"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

// WordLearningHistoryController defines the interface for endpoints that handle
// word learning history records, including tracking quiz results and user progress
type WordLearningHistoryController interface {
	// ProcessQuizResults handles POST requests to process and store quiz results for a user
	ProcessQuizResults(c *gin.Context)
}

// WordLearningHistoryControllerImpl implements the WordLearningHistoryController interface
// and handles operations related to tracking and managing word learning progress
type WordLearningHistoryControllerImpl struct {
	// Service is the business logic layer for word learning history operations
	Service services.WordLearningHistoryService
}

// Make sure that WordLearningHistoryControllerImpl implements WordLearningHistoryController
var _ WordLearningHistoryController = (*WordLearningHistoryControllerImpl)(nil)

// ProcessQuizResults godoc
// @Summary Process quiz results for a user
// @Description Updates learning history for multiple words based on quiz results
// @Tags app
// @Accept json
// @Produce json
// @Param quizResults body dto.QuizResults true "Quiz results"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/app/quiz/results [post]
//
// ProcessQuizResults handles POST requests to submit and process quiz results.
// It extracts the user ID from the authentication token and updates the word learning history
// based on the submitted quiz results.
//
// The request body must contain a QuizResults JSON structure with individual word results.
// On success, it returns HTTP 200 OK with no content.
//
// Possible responses:
//   - 200 OK: Quiz results successfully processed
//   - 400 Bad Request: Invalid request format
//   - 401 Unauthorized: Missing or invalid authentication token
//   - 500 Internal Server Error: Error processing quiz results
func (ctrl *WordLearningHistoryControllerImpl) ProcessQuizResults(c *gin.Context) {
	var quizResults dto.QuizResults
	if err := c.ShouldBindJSON(&quizResults); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user ID from token
	userID, err := middlewares.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to get user ID from token"})
		return
	}

	if err := ctrl.Service.ProcessQuizResults(userID, quizResults.Results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
