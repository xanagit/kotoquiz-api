package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/xanagit/kotoquiz-api/models"
	"net/http"
	"testing"
)

func TestCreateTag(t *testing.T) {
	t.Parallel()

	var label models.Label
	httpResCode := post("/api/v1/tags", `{"en": "Test Tag", "fr": "Tag de Test", "type": "TAG"}`, &label)

	assert.Equal(t, http.StatusCreated, httpResCode)
	assert.Equal(t, "Test Tag", label.En)
	assert.Equal(t, "Tag de Test", label.Fr)
	assert.Equal(t, "TAG", label.Type)
}

func TestReadTag(t *testing.T) {
	t.Parallel()
	var httpResCode int

	var label models.Label
	httpResCode = post("/api/v1/tags", `{"en": "Test Tag", "fr": "Tag de Test", "type": "TAG"}`, &label)

	assert.Equal(t, http.StatusCreated, httpResCode)
	var readLabel models.Label
	httpResCode = get("/api/v1/tags/"+label.ID.String(), &readLabel)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.Equal(t, label.ID, readLabel.ID)
	assert.Equal(t, "Test Tag", readLabel.En)
	assert.Equal(t, "Tag de Test", readLabel.Fr)
	assert.Equal(t, "TAG", readLabel.Type)
}
