package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/xanagit/kotoquiz-api/models"
	"net/http"
	"testing"
)

func TestCreateLabel(t *testing.T) {
	t.Parallel()

	var label models.Label
	httpResCode := post("/api/v1/tech/labels", `{"en": "Test Label", "fr": "Label de Test", "type": "TAG"}`, &label)

	assert.Equal(t, http.StatusCreated, httpResCode)
	assert.Equal(t, "Test Label", label.En)
	assert.Equal(t, "Label de Test", label.Fr)
	assert.Equal(t, "TAG", label.Type)
}

func TestReadLabel(t *testing.T) {
	t.Parallel()
	var httpResCode int

	var label models.Label
	httpResCode = post("/api/v1/tech/labels", `{"en": "Test Label", "fr": "Label de Test", "type": ""}`, &label)

	assert.Equal(t, http.StatusCreated, httpResCode)
	var readLabel models.Label
	httpResCode = get("/api/v1/labels/"+label.ID.String(), &readLabel)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.Equal(t, label.ID, readLabel.ID)
	assert.Equal(t, "Test Label", readLabel.En)
	assert.Equal(t, "Label de Test", readLabel.Fr)
	assert.Equal(t, "TAG", readLabel.Type)
}

func TestListLabel(t *testing.T) {
	t.Parallel()
	var httpResCode int

	initLabels := make([]models.Label, 3)
	httpResCode = post("/api/v1/labels", `{"en": "Test Label1", "fr": "Label de Test1", "type": ""}`, &initLabels[0])
	httpResCode = post("/api/v1/labels", `{"en": "Test Label2", "fr": "Label de Test2", "type": ""}`, &initLabels[1])
	httpResCode = post("/api/v1/labels", `{"en": "Test Label3", "fr": "Label de Test3", "type": ""}`, &initLabels[2])

	var fetchedLabels []models.Label
	httpResCode = get("/api/v1/labels", &fetchedLabels)
	assert.Equal(t, http.StatusOK, httpResCode)

	for _, label := range initLabels {
		assertLabelExistsInList(t, label, fetchedLabels)
	}
}

func TestUpdateLabel(t *testing.T) {
	t.Parallel()

	var label models.Label
	post("/api/v1/labels", `{"en": "Test Label", "fr": "Label de Test", "type": ""}`, &label)

	var updatedLabel models.Label
	httpResCode := put("/api/v1/labels/"+label.ID.String(), `{"id": "99999999-0000-0000-0000-000000000000", "en": "Updated Test Label", "fr": "Label de Test Modifié", "type": ""}`, &updatedLabel)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.Equal(t, label.ID, updatedLabel.ID)
	assert.Equal(t, "Updated Test Label", updatedLabel.En)
	assert.Equal(t, "Label de Test Modifié", updatedLabel.Fr)
	assert.Equal(t, "TAG", updatedLabel.Type)
}

func TestDeleteLabel(t *testing.T) {
	t.Parallel()

	var label models.Label
	post("/api/v1/labels", `{"en": "Test Label", "fr": "Label de Test", "type": ""}`, &label)

	httpResCode := del("/api/v1/labels/" + label.ID.String())

	assert.Equal(t, http.StatusNoContent, httpResCode)
}

func assertLabelExistsInList(t *testing.T, label models.Label, labels []models.Label) {
	for _, currLabel := range labels {
		if currLabel.ID == label.ID {
			assert.Equal(t, label.ID, currLabel.ID)
			assert.Equal(t, label.En, currLabel.En)
			assert.Equal(t, label.Fr, currLabel.Fr)
			assert.Equal(t, label.Type, currLabel.Type)
		}
	}
}
