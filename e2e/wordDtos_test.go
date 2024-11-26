package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/models"
	"net/http"
	"strconv"
	"testing"
)

func Test_should_list_word_ids(t *testing.T) {
	t.Parallel()

	word := GenerateWord()
	var insertedWord models.Word
	httpResCode := post("/api/v1/tech/words", ToJson(&word), &insertedWord)
	assert.Equal(t, http.StatusCreated, httpResCode)

	var wordsIds dto.WordIdsList
	httpResCode = get("/api/v1/app/words/"+insertedWord.ID.String()+"?lang=fr", &wordsIds)
}

func Test_should_read_wordDto(t *testing.T) {
	t.Parallel()

	word := GenerateWord()
	var insertedWord models.Word
	httpResCode := post("/api/v1/tech/words", ToJson(&word), &insertedWord)
	assert.Equal(t, http.StatusCreated, httpResCode)

	var fetchedWordDto dto.WordDTO
	httpResCode = get("/api/v1/app/words/"+insertedWord.ID.String()+"?lang=fr", &fetchedWordDto)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.Equal(t, insertedWord.ID, fetchedWordDto.ID)
	assert.Equal(t, insertedWord.Kanji, fetchedWordDto.Kanji)
	assert.Equal(t, insertedWord.Yomi, fetchedWordDto.Yomi)
	assert.Equal(t, insertedWord.YomiType, fetchedWordDto.YomiType)
	assert.Equal(t, insertedWord.ImageURL, fetchedWordDto.ImageURL)
	assert.Equal(t, insertedWord.Translation.Fr, fetchedWordDto.Translation)
	assert.Equal(t, len(insertedWord.Tags), len(fetchedWordDto.Tags))
	for idx, tag := range insertedWord.Tags {
		assert.Equal(t, tag.Fr, fetchedWordDto.Tags[idx])
	}
	assert.Equal(t, len(insertedWord.Levels), len(fetchedWordDto.Levels))
	for idx, level := range insertedWord.Levels {
		assert.Equal(t, level.Category.Fr, fetchedWordDto.Levels[idx].Category)
		assert.Equal(t, len(level.LevelNames), len(fetchedWordDto.Levels[idx].LevelNames))
		for idx, levelName := range level.LevelNames {
			assert.Equal(t, levelName.Fr, fetchedWordDto.Levels[idx].LevelNames[idx])
		}
	}
}

func Test_should_list_WordDtoIds_corresponding_to_provided_tag_ids(t *testing.T) {
	t.Parallel()

	insertedWords := insertWordsDatasetForListDtoIds(t)

	var fetchedWordDtoIdsList dto.WordIdsList
	httpResCode := get("/api/v1/app/words/q?lang=fr&tags="+insertedWords[0].Tags[0].ID.String(), &fetchedWordDtoIdsList)
	assert.Equal(t, http.StatusOK, httpResCode)

	// Check that words only corresponding to tag are fetched
	assert.Equal(t, 2, len(fetchedWordDtoIdsList.Ids))
	assert.Contains(t, fetchedWordDtoIdsList.Ids, insertedWords[0].ID.String())
	assert.Contains(t, fetchedWordDtoIdsList.Ids, insertedWords[2].ID.String())
	assert.NotContains(t, fetchedWordDtoIdsList.Ids, insertedWords[1].ID.String())
}

func Test_should_list_WordDtoIds_corresponding_to_provided_levelNamesIds_ids(t *testing.T) {
	t.Parallel()

	insertedWords := insertWordsDatasetForListDtoIds(t)

	var fetchedWordDtoIdsForLevelsList dto.WordIdsList
	httpResCode := get("/api/v1/app/words/q?lang=fr&levelNames="+insertedWords[0].Levels[0].LevelNames[0].ID.String(), &fetchedWordDtoIdsForLevelsList)
	assert.Equal(t, http.StatusOK, httpResCode)

	// Check that words only corresponding to levelNames are fetched
	assert.Equal(t, 2, len(fetchedWordDtoIdsForLevelsList.Ids))
	assert.Contains(t, fetchedWordDtoIdsForLevelsList.Ids, insertedWords[0].ID.String())
	assert.Contains(t, fetchedWordDtoIdsForLevelsList.Ids, insertedWords[1].ID.String())
	assert.NotContains(t, fetchedWordDtoIdsForLevelsList.Ids, insertedWords[2].ID.String())
}

func Test_should_list_WordDtos(t *testing.T) {
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

	wIds := insertedWords[0].ID.String() + "," + insertedWords[1].ID.String() + "," + insertedWords[2].ID.String()
	var fetchedWordDtos []dto.WordDTO
	httpResCode = get("/api/v1/app/words?lang=fr&ids="+wIds, &fetchedWordDtos)
	assert.Equal(t, http.StatusOK, httpResCode)

	for _, word := range insertedWords {
		assertWordExistsInWordDtoList(t, word, fetchedWordDtos)
	}
}

func assertWordExistsInWordDtoList(t *testing.T, word models.Word, wordDtos []dto.WordDTO) {
	for _, currWordDto := range wordDtos {
		if currWordDto.ID == word.ID {
			assert.Equal(t, word.ID, currWordDto.ID)
			assert.Equal(t, word.Kanji, currWordDto.Kanji)
			assert.Equal(t, word.Yomi, currWordDto.Yomi)
			assert.Equal(t, word.YomiType, currWordDto.YomiType)
			assert.Equal(t, word.ImageURL, currWordDto.ImageURL)
			assert.Equal(t, word.Translation.Fr, currWordDto.Translation)
			assert.Equal(t, len(word.Tags), len(currWordDto.Tags))
			for idx, tag := range word.Tags {
				assert.Equal(t, tag.Fr, currWordDto.Tags[idx])
			}
			assert.Equal(t, len(word.Levels), len(currWordDto.Levels))
			for idx, level := range word.Levels {
				assert.Equal(t, level.Category.Fr, currWordDto.Levels[idx].Category)
				assert.Equal(t, len(level.LevelNames), len(currWordDto.Levels[idx].LevelNames))
				for idx, levelName := range level.LevelNames {
					assert.Equal(t, levelName.Fr, currWordDto.Levels[idx].LevelNames[idx])
				}
			}
		}
	}
}

func insertWordsDatasetForListDtoIds(t *testing.T) []*models.Word {
	var words []*models.Word
	for i := 0; i < 3; i++ {
		word := GenerateWord()
		words = append(words, &word)
	}

	tag := generateLabel(models.Tag)
	var insertedTag models.Label
	httpResCode := post("/api/v1/tech/tags", ToJson(&tag), &insertedTag)
	assert.Equal(t, http.StatusCreated, httpResCode)

	level := GenerateLevel()
	var insertedLevel models.Level
	httpResCode = post("/api/v1/tech/levels", ToJson(&level), &insertedLevel)
	assert.Equal(t, http.StatusCreated, httpResCode)

	for idx, word := range words {
		word.Translation.En = "En" + strconv.Itoa(idx)
		word.Translation.Fr = "Fr" + strconv.Itoa(idx)
		word.Tags = nil
		word.Levels = nil
		if idx != 1 { // No tag for second word
			word.Tags = []*models.Label{&insertedTag}
		}
		if idx != 2 { // No level for third word
			word.Levels = []*models.Level{&insertedLevel}
		}
	}

	insertedWords := make([]*models.Word, 3)
	for idx, word := range words {
		httpResCode = post("/api/v1/tech/words", ToJson(&word), &insertedWords[idx])
		assert.Equal(t, http.StatusCreated, httpResCode)
	}

	return insertedWords
}
