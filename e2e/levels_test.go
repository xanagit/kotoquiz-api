package main

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xanagit/kotoquiz-api/models"
	"net/http"
	"strconv"
	"testing"
)

func Test_should_create_level(t *testing.T) {
	t.Parallel()

	level := GenerateLevel()
	var resLevel models.Level
	httpResCode := post("/api/v1/tech/levels", ToJson(&level), &resLevel)

	assert.Equal(t, http.StatusCreated, httpResCode)
	assert.NotEqual(t, level.ID, resLevel.ID)
	level.ID = resLevel.ID
	assert.Equal(t, level, resLevel)
}

func Test_should_read_level(t *testing.T) {
	t.Parallel()

	level := GenerateLevel()
	var insertedLevel models.Level
	httpResCode := post("/api/v1/tech/levels", ToJson(&level), &insertedLevel)

	assert.Equal(t, http.StatusCreated, httpResCode)
	var fetchedLevel models.Level
	httpResCode = get("/api/v1/tech/levels/"+insertedLevel.ID.String(), &fetchedLevel)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.Equal(t, insertedLevel.ID, fetchedLevel.ID)
	assert.Equal(t, insertedLevel, fetchedLevel)
}

func Test_should_update_level(t *testing.T) {
	t.Parallel()

	level := GenerateLevel()
	level.ID = uuid.Nil
	var insertedTag models.Level
	post("/api/v1/tech/levels", ToJson(&level), &insertedTag)

	var updatedTag models.Level
	level.ID, _ = uuid.Parse("99999999-9999-9999-9999-999999999999")
	level.Category.En = "En Updated"
	level.Category.Fr = "Fr Modifié"
	level.LevelNames[0].En = "LevelNames En Updated"
	level.LevelNames[0].Fr = "LevelNames Fr Modifié"
	level.LevelNames[1].En = "LevelNames En Updated"
	level.LevelNames[1].Fr = "LevelNames Fr Modifié"

	httpResCode := put("/api/v1/tech/levels/"+insertedTag.ID.String(), ToJson(&level), &updatedTag)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.NotEqual(t, level.ID, updatedTag.ID)
	assert.Equal(t, insertedTag.ID, updatedTag.ID)
	level.ID = updatedTag.ID
	assert.Equal(t, level, updatedTag)
}

func Test_should_delete_level(t *testing.T) {
	t.Parallel()

	level := GenerateLevel()
	level.ID = uuid.Nil
	var insertedLevel models.Level
	post("/api/v1/tech/levels", ToJson(&level), &insertedLevel)

	httpResCode := del("/api/v1/tech/levels/" + insertedLevel.ID.String())

	assert.Equal(t, http.StatusNoContent, httpResCode)

	httpResCode = get("/api/v1/tech/levels/"+insertedLevel.ID.String(), &insertedLevel)
	assert.Equal(t, http.StatusNotFound, httpResCode)
}

func Test_should_list_Levels(t *testing.T) {
	t.Parallel()
	var httpResCode int

	levels := []models.Level{GenerateLevel(), GenerateLevel(), GenerateLevel()}

	insertedLevels := make([]models.Level, 3)
	for idx, level := range levels {
		level.Category.En = "En" + strconv.Itoa(idx)
		level.Category.Fr = "Fr" + strconv.Itoa(idx)
		for idx, l := range level.LevelNames {
			l.En = "LevelNames En " + strconv.Itoa(idx)
			l.Fr = "LevelNames Fr " + strconv.Itoa(idx)
		}
		httpResCode = post("/api/v1/tech/levels", ToJson(&level), &insertedLevels[idx])
		assert.Equal(t, http.StatusCreated, httpResCode)
	}

	var fetchedTags []models.Level
	httpResCode = get("/api/v1/app/levels", &fetchedTags)
	assert.Equal(t, http.StatusOK, httpResCode)

	for _, level := range insertedLevels {
		assertLevelExistsInList(t, level, fetchedTags)
	}
}

func assertLevelExistsInList(t *testing.T, level models.Level, levels []models.Level) {
	for _, currLevel := range levels {
		if currLevel.ID == level.ID {
			assert.Equal(t, level, currLevel)
		}
	}
}

func GenerateLevel() models.Level {
	restInputLevel := models.Level{
		ID: uuid.New(),
		Category: models.Label{
			ID:   uuid.New(),
			En:   "Category En",
			Fr:   "Category Fr",
			Type: models.Category,
		},
		LevelNames: []*models.Label{
			{
				ID:   uuid.New(),
				En:   "LevelNames En 1",
				Fr:   "LevelNames Fr 1",
				Type: models.LevelName,
			}, {
				ID:   uuid.New(),
				En:   "LevelNames En 2",
				Fr:   "LevelNames Fr 2",
				Type: models.LevelName,
			}},
	}
	return restInputLevel
}
