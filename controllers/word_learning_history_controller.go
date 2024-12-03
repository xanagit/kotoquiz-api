package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type WordLearningHistoryController interface {
	ProcessQuizResults(c *gin.Context)
}

type WordLearningHistoryControllerImpl struct {
	Service services.WordLearningHistoryService
}

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
func (ctrl *WordLearningHistoryControllerImpl) ProcessQuizResults(c *gin.Context) {
	var quizResults dto.QuizResults
	if err := c.ShouldBindJSON(&quizResults); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Service.ProcessQuizResults(quizResults.UserID, quizResults.Results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
