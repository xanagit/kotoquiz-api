package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type defaultValues struct {
	Lang        string
	NbIdsList   int
	LimitWords  int
	OffsetWords int
}

const (
	DefaultLang        = "en"
	DefaultNbIdsList   = 30
	DefaultLimitWords  = 15
	DefaultOffsetWords = 0
)

var DefaultQpVals = defaultValues{
	Lang:        DefaultLang,
	NbIdsList:   DefaultNbIdsList,
	LimitWords:  DefaultLimitWords,
	OffsetWords: DefaultOffsetWords,
}

func FromStrToUuid(id string) uuid.UUID {
	parsed, err := uuid.Parse(id)
	if err != nil {
		log.Fatalf("Invalid UUID format: %v", err)
	}
	return parsed
}

func getQueryParamList(c *gin.Context, paramName string) []string {
	rawList := c.Query(paramName)
	var strList []string
	if rawList != "" {
		strList = strings.Split(rawList, ",")
	}
	return strList
}

func getQueryParamInt(c *gin.Context, paramName string, defaultValue int) (int, error) {
	rawParam := c.Param(paramName)
	param := defaultValue

	var err error
	if rawParam != "" {
		param, err = strconv.Atoi(rawParam)
	}

	if err != nil || param < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid '" + paramName + "' parameter"})
		return 0, err
	}
	return param, nil
}

func getQueryParamLang(c *gin.Context) string {
	lang := c.Query("lang")
	if lang == "" {
		lang = DefaultQpVals.Lang
	}
	return lang
}
