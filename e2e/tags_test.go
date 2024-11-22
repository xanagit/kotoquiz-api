package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/xanagit/kotoquiz-api/models"
	"net/http"
	"testing"
)

func TestCreateTag(t *testing.T) {
	t.Parallel()

	var tag models.Label
	httpResCode := post("/api/v1/tech/tags", `{"en": "Test Tag", "fr": "Tag de Test", "type": ""}`, &tag)

	assert.Equal(t, http.StatusCreated, httpResCode)
	assert.Equal(t, "Test Tag", tag.En)
	assert.Equal(t, "Tag de Test", tag.Fr)
	assert.Equal(t, "TAG", tag.Type)
}

func TestReadTag(t *testing.T) {
	t.Parallel()
	var httpResCode int

	var insertedTag models.Label
	httpResCode = post("/api/v1/tech/tags", `{"en": "Test Tag", "fr": "Tag de Test", "type": ""}`, &insertedTag)

	assert.Equal(t, http.StatusCreated, httpResCode)
	var readLabel models.Label
	httpResCode = get("/api/v1/tech/tags/"+insertedTag.ID.String(), &readLabel)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.Equal(t, insertedTag.ID, readLabel.ID)
	assert.Equal(t, "Test Tag", readLabel.En)
	assert.Equal(t, "Tag de Test", readLabel.Fr)
	assert.Equal(t, "TAG", readLabel.Type)
}

func TestListTag(t *testing.T) {
	t.Parallel()
	var httpResCode int

	insertedTags := make([]models.Label, 3)
	httpResCode = post("/api/v1/tech/tags", `{"en": "Test Tag1", "fr": "Tag de Test1", "type": ""}`, &insertedTags[0])
	httpResCode = post("/api/v1/tech/tags", `{"en": "Test Tag2", "fr": "Tag de Test2", "type": ""}`, &insertedTags[1])
	httpResCode = post("/api/v1/tech/tags", `{"en": "Test Tag3", "fr": "Tag de Test3", "type": ""}`, &insertedTags[2])

	var fetchedTags []models.Label
	httpResCode = get("/api/v1/app/tags", &fetchedTags)
	assert.Equal(t, http.StatusOK, httpResCode)

	for _, tag := range insertedTags {
		assertTagExistsInList(t, tag, fetchedTags)
	}
}

func TestUpdateTag(t *testing.T) {
	t.Parallel()

	var insertedTag models.Label
	post("/api/v1/tech/tags", `{"en": "Test Tag", "fr": "Tag de Test", "type": ""}`, &insertedTag)

	var updatedLabel models.Label
	httpResCode := put("/api/v1/tech/tags/"+insertedTag.ID.String(), `{"id": "99999999-0000-0000-0000-000000000000", "en": "Updated Test Tag", "fr": "Tag de Test Modifié", "type": ""}`, &updatedLabel)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.Equal(t, insertedTag.ID, updatedLabel.ID)
	assert.Equal(t, "Updated Test Tag", updatedLabel.En)
	assert.Equal(t, "Tag de Test Modifié", updatedLabel.Fr)
	assert.Equal(t, "TAG", updatedLabel.Type)
}

func TestDeleteTag(t *testing.T) {
	t.Parallel()

	var label models.Label
	post("/api/v1/tech/tags", `{"en": "Test Tag", "fr": "Tag de Test", "type": ""}`, &label)

	httpResCode := del("/api/v1/tech/tags/" + label.ID.String())

	assert.Equal(t, http.StatusNoContent, httpResCode)
}

func assertTagExistsInList(t *testing.T, tag models.Label, tags []models.Label) {
	for _, currTag := range tags {
		if currTag.ID == tag.ID {
			assert.Equal(t, tag.ID, currTag.ID)
			assert.Equal(t, tag.En, currTag.En)
			assert.Equal(t, tag.Fr, currTag.Fr)
			assert.Equal(t, tag.Type, currTag.Type)
		}
	}
}
