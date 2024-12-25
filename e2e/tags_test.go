package main

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xanagit/kotoquiz-api/models"
	"net/http"
	"strconv"
	"testing"
)

func Test_should_create_tag(t *testing.T) {
	t.Parallel()

	tag := generateLabel(models.Tag)
	tag.Type = models.Translation // To test API auto-assigning the correct type
	var insertedTag models.Label
	httpResCode := post("/api/v1/tech/tags", ToJson(&tag), &insertedTag)

	assert.Equal(t, http.StatusCreated, httpResCode)
	assert.NotEqual(t, tag.ID, insertedTag.ID)
	tag.ID = insertedTag.ID
	tag.Type = models.Tag
	assert.Equal(t, tag, insertedTag)
}

func Test_should_read_tag(t *testing.T) {
	t.Parallel()
	var httpResCode int

	tag := generateLabel(models.Tag)
	tag.ID = uuid.Nil
	tag.Type = models.Translation // To test API auto-assigning the correct type
	var insertedTag models.Label
	httpResCode = post("/api/v1/tech/tags", ToJson(&tag), &insertedTag)

	assert.Equal(t, http.StatusCreated, httpResCode)
	var fetchedTag models.Label
	httpResCode = get("/api/v1/tech/tags/"+insertedTag.ID.String(), &fetchedTag)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.Equal(t, insertedTag, fetchedTag)
	assert.Equal(t, models.Tag, fetchedTag.Type)
}

func Test_should_update_tag(t *testing.T) {
	t.Parallel()

	tag := generateLabel(models.Tag)
	tag.ID = uuid.Nil
	var insertedTag models.Label
	post("/api/v1/tech/tags", ToJson(&tag), &insertedTag)

	var updatedTag models.Label
	tag.ID = uuid.New()
	tag.En = "En Updated"
	tag.Fr = "Fr Modifié"
	httpResCode := put("/api/v1/tech/tags/"+insertedTag.ID.String(), ToJson(&tag), &updatedTag)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.NotEqual(t, tag.ID, updatedTag.ID)
	assert.Equal(t, insertedTag.ID, updatedTag.ID)
	tag.ID = updatedTag.ID
	assert.Equal(t, tag, updatedTag)
	assert.Equal(t, models.Tag, updatedTag.Type)
}

func Test_should_delete_tag(t *testing.T) {
	t.Parallel()

	tag := generateLabel(models.Tag)
	tag.ID = uuid.Nil
	var insertedTag models.Label
	post("/api/v1/tech/tags", ToJson(&tag), &insertedTag)

	httpResCode := del("/api/v1/tech/tags/" + insertedTag.ID.String())

	assert.Equal(t, http.StatusNoContent, httpResCode)
}

func Test_should_list_Tags(t *testing.T) {
	t.Parallel()
	var httpResCode int

	tags := []models.Label{generateLabel(models.Tag), generateLabel(models.Tag), generateLabel(models.Tag)}

	insertedTags := make([]models.Label, 3)
	for idx, tag := range tags {
		tag.En = "En" + strconv.Itoa(idx)
		tag.Fr = "Fr" + strconv.Itoa(idx)
		httpResCode = post("/api/v1/tech/tags", ToJson(&tag), &insertedTags[idx])
		assert.Equal(t, http.StatusCreated, httpResCode)
	}

	var fetchedTags []models.Label
	httpResCode = get("/api/v1/app/tags", &fetchedTags)
	assert.Equal(t, http.StatusOK, httpResCode)

	for _, tag := range insertedTags {
		assertTagExistsInList(t, tag, fetchedTags)
	}
}

func assertTagExistsInList(t *testing.T, tag models.Label, tags []models.Label) {
	for _, currTag := range tags {
		if currTag.ID == tag.ID {
			assert.Equal(t, tag, currTag)
		}
	}
}

func generateLabel(labelType models.LabelType) models.Label {
	restInputLabel := models.Label{
		ID:   uuid.New(),
		En:   "Label En",
		Fr:   "Label Fr",
		Type: labelType,
	}
	return restInputLabel
}
