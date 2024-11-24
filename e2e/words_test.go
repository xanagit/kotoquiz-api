package main

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xanagit/kotoquiz-api/models"
	"net/http"
	"strconv"
	"testing"
)

func Test_should_create_word(t *testing.T) {
	t.Parallel()

	word := GenerateWord()
	word.Translation.Type = ""
	var resWord models.Word
	httpResCode := post("/api/v1/tech/words", ToJson(&word), &resWord)

	assert.Equal(t, http.StatusCreated, httpResCode)
	assert.NotEqual(t, word.ID, resWord.ID)
	word.ID = resWord.ID
	word.Translation.Type = models.Translation
	assert.Equal(t, word, resWord)
}

func Test_should_read_word(t *testing.T) {
	t.Parallel()

	word := GenerateWord()
	var insertedWord models.Word
	httpResCode := post("/api/v1/tech/words", ToJson(&word), &insertedWord)

	assert.Equal(t, http.StatusCreated, httpResCode)
	var fetchedWord models.Word
	httpResCode = get("/api/v1/tech/words/"+insertedWord.ID.String(), &fetchedWord)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.Equal(t, insertedWord.ID, fetchedWord.ID)
	assert.Equal(t, insertedWord, fetchedWord)
}

func Test_should_update_word(t *testing.T) {
	t.Parallel()

	word := GenerateWord()
	word.ID = uuid.Nil
	var insertedTag models.Word
	post("/api/v1/tech/words", ToJson(&word), &insertedTag)

	var updatedTag models.Word
	word.ID, _ = uuid.Parse("99999999-9999-9999-9999-999999999999")
	word.Kanji = "Kanji Updated"
	word.Yomi = "Yomi Updated"
	word.YomiType = models.Kunyomi
	word.ImageURL = "https://kotoquiz.com/image_updated.jpg"
	word.Translation.En = "Translation En Updated"
	word.Translation.Fr = "Translation Fr Updated"
	word.Tags[0].En = "Tag En Updated 1"
	word.Tags[0].Fr = "Tag Fr Updated 1"
	word.Tags[1].En = "Tag En Updated 2"
	word.Tags[1].Fr = "Tag Fr Updated 2"
	word.Levels[0].Category.En = "Category En Updated 1"
	word.Levels[0].Category.Fr = "Category Fr Updated 1"
	word.Levels[0].LevelNames[0].En = "LevelNames En Updated 1"
	word.Levels[0].LevelNames[0].Fr = "LevelNames Fr Updated 1"
	word.Levels[0].LevelNames[1].En = "LevelNames En Updated 2"
	word.Levels[0].LevelNames[1].Fr = "LevelNames Fr Updated 2"

	httpResCode := put("/api/v1/tech/words/"+insertedTag.ID.String(), ToJson(&word), &updatedTag)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.NotEqual(t, word.ID, updatedTag.ID)
	assert.Equal(t, insertedTag.ID, updatedTag.ID)
	word.ID = updatedTag.ID
	assert.Equal(t, word, updatedTag)
}

func Test_should_delete_word(t *testing.T) {
	t.Parallel()

	word := GenerateWord()
	word.ID = uuid.Nil
	var insertedWord models.Word
	post("/api/v1/tech/words", ToJson(&word), &insertedWord)

	httpResCode := del("/api/v1/tech/words/" + insertedWord.ID.String())

	assert.Equal(t, http.StatusNoContent, httpResCode)

	httpResCode = get("/api/v1/tech/words/"+insertedWord.ID.String(), &insertedWord)
	assert.Equal(t, http.StatusNotFound, httpResCode)

	// Checking if levels still exist
	wordLevels := []models.Level{*word.Levels[0], *word.Levels[1]}
	//var levels []models.Level
	remainingLevels := make([]models.Level, 3)
	httpResCode = get("/api/v1/app/levels", &remainingLevels)
	assert.Equal(t, http.StatusOK, httpResCode)

	for _, level := range remainingLevels {
		assertLevelExistsInList(t, level, wordLevels)
	}

	// Checking if tags still exist
	wordTag := []models.Label{*word.Tags[0], *word.Tags[1]}
	remainingTags := make([]models.Label, 3)
	httpResCode = get("/api/v1/app/tags", &remainingTags)
	assert.Equal(t, http.StatusOK, httpResCode)

	for _, tag := range remainingTags {
		assertTagExistsInList(t, tag, wordTag)
	}
}

func Test_should_list_Words(t *testing.T) {
	t.Parallel()
	var httpResCode int

	words := []models.Word{GenerateWord(), GenerateWord(), GenerateWord()}

	insertedWords := make([]models.Word, 3)
	for idx, word := range words {
		word.Translation.En = "En" + strconv.Itoa(idx)
		word.Translation.Fr = "Fr" + strconv.Itoa(idx)
		for idx, t := range word.Tags {
			t.En = "Tag En " + strconv.Itoa(idx)
			t.Fr = "Tag Fr " + strconv.Itoa(idx)
		}
		for idx, l := range word.Levels {
			l.Category.En = "Category En " + strconv.Itoa(idx)
			l.Category.Fr = "Category Fr " + strconv.Itoa(idx)
			for idx, ln := range l.LevelNames {
				ln.En = "LevelNames En " + strconv.Itoa(idx)
				ln.Fr = "LevelNames Fr " + strconv.Itoa(idx)
			}
		}
		httpResCode = post("/api/v1/tech/words", ToJson(&word), &insertedWords[idx])
		assert.Equal(t, http.StatusCreated, httpResCode)
	}

	var fetchedWords []models.Word
	httpResCode = get("/api/v1/app/words", &fetchedWords)
	assert.Equal(t, http.StatusOK, httpResCode)

	for _, word := range insertedWords {
		assertWordExistsInList(t, word, fetchedWords)
	}
}

func assertWordExistsInList(t *testing.T, word models.Word, words []models.Word) {
	for _, currWord := range words {
		if currWord.ID == word.ID {
			assert.Equal(t, word, currWord)
		}
	}
}

func GenerateWord() models.Word {
	levels := []models.Level{GenerateLevel(), GenerateLevel(), GenerateLevel()}
	restInputWord := models.Word{
		ID:       uuid.New(),
		Kanji:    "kanki",
		Yomi:     "yomi",
		YomiType: models.Onyomi,
		ImageURL: "https://kotoquiz.com/image.jpg",
		Translation: models.Label{
			ID:   uuid.New(),
			En:   "Translation En",
			Fr:   "Translation Fr",
			Type: models.Translation,
		},
		Tags: []*models.Label{
			{
				ID:   uuid.New(),
				En:   "Tag En 1",
				Fr:   "Tag Fr 1",
				Type: models.Tag,
			}, {
				ID:   uuid.New(),
				En:   "Tag En 2",
				Fr:   "Tag Fr 2",
				Type: models.Tag,
			},
		},
		Levels: []*models.Level{
			&levels[0], &levels[1], &levels[2],
		},
	}
	GenerateLevel()
	return restInputWord
}
