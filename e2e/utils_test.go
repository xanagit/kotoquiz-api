package main

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
)

func get[T any](url string, model *T) int {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	router.ServeHTTP(w, req)

	err := json.Unmarshal(w.Body.Bytes(), model)
	if err != nil {
		logger.Error("Could not unmarshall json")
	}

	return w.Code
}

func post[T any](url string, jsonData string, model *T) int {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(jsonData))
	router.ServeHTTP(w, req)

	err := json.Unmarshal(w.Body.Bytes(), model)
	if err != nil {
		logger.Error("Could not unmarshall json")
	}
	logger.Info("label", zap.Any("model", model))

	return w.Code
}

func put[T any](url string, jsonData string, model *T) int {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, bytes.NewBufferString(jsonData))
	router.ServeHTTP(w, req)

	err := json.Unmarshal(w.Body.Bytes(), model)
	if err != nil {
		logger.Error("Could not unmarshall json")
	}
	logger.Info("label", zap.Any("model", model))

	return w.Code
}

func del(url string) int {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", url, nil)
	router.ServeHTTP(w, req)

	return w.Code
}

func ToJson[T any](input *T) string {
	jsonData, _ := json.MarshalIndent(input, "", "  ")
	return string(jsonData)
}
