package dto

import "github.com/google/uuid"

type ResultStatus string

const (
	Success    ResultStatus = "SUCCESS"
	Error      ResultStatus = "ERROR"
	Unanswered ResultStatus = "UNANSWERED"
)

type WordQuizResult struct {
	WordID uuid.UUID    `json:"wordId"`
	Status ResultStatus `json:"type"`
}

type QuizResults struct {
	UserID  uuid.UUID        `json:"userId"`
	Results []WordQuizResult `json:"results"`
}
