package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"strings"
)

func FromStrToUuid(id string) uuid.UUID {
	parsed, err := uuid.Parse(id)
	if err != nil {
		log.Fatalf("Invalid UUID format: %v", err)
	}
	return parsed
}

func getQueryParamIds(c *gin.Context) []string {
	idsParam := c.Query("ids")
	var ids []string
	if idsParam != "" {
		ids = strings.Split(idsParam, ",")
	}
	return ids
}

func getQueryParamLang(c *gin.Context) string {
	lang := c.Query("lang")
	if lang == "" {
		lang = "en"
	}
	return lang
}
