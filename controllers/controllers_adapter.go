// Package controllers implements the HTTP handlers for the application API endpoints
// It defines interfaces and implementations for all API controllers, handling request
// parsing, validation, and response formatting.
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// getQueryParamList extracts a comma-separated query parameter as a string slice
//
// Parameters:
//   - c: *gin.Context - The Gin context containing the request
//   - key: string - The name of the query parameter to extract
//
// Returns:
//   - []string - A slice of strings from the comma-separated parameter, or empty slice if not found
func getQueryParamList(c *gin.Context, paramName string) []string {
	rawList := c.Query(paramName)
	var strList []string
	if rawList != "" {
		strList = strings.Split(rawList, ",")
	}
	return strList
}

// getQueryParamInt extracts an integer query parameter with a default value
//
// Parameters:
//   - c: *gin.Context - The Gin context containing the request
//   - key: string - The name of the query parameter to extract
//   - defaultValue: int - The default value to return if parameter is missing or invalid
//
// Returns:
//   - int - The parsed integer value or default if not found/invalid
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

// getQueryParamLang extracts a language code from query parameters, defaulting to Japanese
//
// Parameters:
//   - c: *gin.Context - The Gin context containing the request
//
// Returns:
//   - string - The language code (defaults to "ja" if not specified)
func getQueryParamLang(c *gin.Context) string {
	lang := c.Query("lang")
	if lang == "" {
		lang = DefaultQpVals.Lang
	}
	return lang
}

// parseUUID parses a string into a UUID, handling errors
//
// Parameters:
//   - uuidStr: string - The string to parse as a UUID
//
// Returns:
//   - uuid.UUID - The parsed UUID
//   - error - An error if parsing fails
func parseUUID(id string) (uuid.UUID, bool) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, false
	}
	return parsed, true
}

// parseUUIDs parses a slice of string UUIDs, handling errors
//
// Parameters:
//   - uuidStrs: []string - A slice of strings to parse as UUIDs
//
// Returns:
//   - []uuid.UUID - A slice of the parsed UUIDs
//   - error - An error if any parsing fails
func parseUUIDs(ids []string) ([]uuid.UUID, bool) {
	parsed := make([]uuid.UUID, len(ids))
	for i, id := range ids {
		currUUID, ok := parseUUID(id)
		if !ok {
			return nil, false
		}
		parsed[i] = currUUID
	}
	return parsed, true
}
