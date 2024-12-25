package main

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/models"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
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

	insertedWords := insertWordsDatasetForListDtoIds(t, 3)
	var fetchedWordDtoIdsList dto.WordIdsList
	httpResCode := get("/api/v1/app/words/q?lang=fr&tags="+insertedWords[0].Tags[0].ID.String()+"&userId="+uuid.New().String(), &fetchedWordDtoIdsList)
	assert.Equal(t, http.StatusOK, httpResCode)

	// Check that words only corresponding to tag are fetched
	assert.Equal(t, 2, len(fetchedWordDtoIdsList.Ids))
	assert.Contains(t, fetchedWordDtoIdsList.Ids, insertedWords[0].ID.String())
	assert.Contains(t, fetchedWordDtoIdsList.Ids, insertedWords[2].ID.String())
	assert.NotContains(t, fetchedWordDtoIdsList.Ids, insertedWords[1].ID.String())
}

func Test_should_list_WordDtoIds_corresponding_to_provided_levelNamesIds_ids(t *testing.T) {
	t.Parallel()

	insertedWords := insertWordsDatasetForListDtoIds(t, 3)

	var fetchedWordDtoIdsForLevelsList dto.WordIdsList
	httpResCode := get("/api/v1/app/words/q?lang=fr&levelNames="+insertedWords[0].Levels[0].LevelNames[0].ID.String()+"&userId="+uuid.New().String(), &fetchedWordDtoIdsForLevelsList)
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

	insertedWords := insertWordsDatasetForListDtoIds(t, 3)

	ids := make([]string, len(insertedWords))
	for i, w := range insertedWords {
		ids[i] = w.ID.String()
	}
	wIds := strings.Join(ids, ",")
	var fetchedWordDtos []dto.WordDTO
	httpResCode = get("/api/v1/app/words?lang=fr&ids="+wIds, &fetchedWordDtos)
	assert.Equal(t, http.StatusOK, httpResCode)

	for _, word := range insertedWords {
		assertWordExistsInWordDtoList(t, word, fetchedWordDtos)
	}
}

func assertWordExistsInWordDtoList(t *testing.T, word *models.Word, wordDtos []dto.WordDTO) {
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
			for idxLevel, level := range word.Levels {
				assert.Equal(t, level.Category.Fr, currWordDto.Levels[idxLevel].Category)
				assert.Equal(t, len(level.LevelNames), len(currWordDto.Levels[idxLevel].LevelNames))
				for idxLevelName, levelName := range level.LevelNames {
					assert.Equal(t, levelName.Fr, currWordDto.Levels[idxLevel].LevelNames[idxLevelName])
				}
			}
		}
	}
}

func insertWordsDatasetForListDtoIds(t *testing.T, nb int) []*models.Word {
	var words []*models.Word
	for i := 0; i < nb; i++ {
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

func Test_should_process_quiz_results(t *testing.T) {
	t.Parallel()

	// Create test words
	words := []models.Word{GenerateWord(), GenerateWord(), GenerateWord()}
	insertedWords := make([]models.Word, len(words))
	for idx, word := range words {
		httpResCode := post("/api/v1/tech/words", ToJson(&word), &insertedWords[idx])
		assert.Equal(t, http.StatusCreated, httpResCode)
	}

	// Create quiz results
	quizResults := dto.QuizResults{
		UserID: uuid.New().String(),
		Results: []dto.WordQuizResult{
			{WordID: insertedWords[0].ID, Status: dto.Success},
			{WordID: insertedWords[1].ID, Status: dto.Error},
			{WordID: insertedWords[2].ID, Status: dto.Unanswered},
		},
	}

	// Submit quiz results
	httpResCode := postNoContent("/api/v1/app/quiz/results", ToJson(&quizResults))
	assert.Equal(t, http.StatusOK, httpResCode)

	// Test successive quiz results for learning progression
	testSuccessiveQuizResults(t, quizResults.UserID, insertedWords[0].ID)
}

func Test_should_handle_invalid_quiz_results(t *testing.T) {
	t.Parallel()

	// Test with invalid user ID
	invalidUserResults := dto.QuizResults{
		UserID: uuid.New().String(), // Non-existent user
		Results: []dto.WordQuizResult{
			{WordID: uuid.New(), Status: dto.Success},
		},
	}
	httpResCode := postNoContent("/api/v1/app/quiz/results", ToJson(&invalidUserResults))
	assert.Equal(t, http.StatusInternalServerError, httpResCode)

	invalidWordResults := dto.QuizResults{
		UserID: uuid.New().String(),
		Results: []dto.WordQuizResult{
			{WordID: uuid.New(), Status: dto.Success}, // Non-existent word
		},
	}
	httpResCode = postNoContent("/api/v1/app/quiz/results", ToJson(&invalidWordResults))
	assert.Equal(t, http.StatusInternalServerError, httpResCode)
}

func testSuccessiveQuizResults(t *testing.T, userID string, wordID uuid.UUID) {
	// Test progression through learning states
	successResults := dto.QuizResults{
		UserID: userID,
		Results: []dto.WordQuizResult{
			{WordID: wordID, Status: dto.Success},
		},
	}

	// Submit multiple successful answers to test progression
	for i := 0; i < 5; i++ {
		httpResCode := postNoContent("/api/v1/app/quiz/results", ToJson(&successResults))
		assert.Equal(t, http.StatusOK, httpResCode)
		time.Sleep(100 * time.Millisecond) // Ensure different timestamps
	}

	// Test error impact
	errorResults := dto.QuizResults{
		UserID: userID,
		Results: []dto.WordQuizResult{
			{WordID: wordID, Status: dto.Error},
		},
	}
	httpResCode := postNoContent("/api/v1/app/quiz/results", ToJson(&errorResults))
	assert.Equal(t, http.StatusOK, httpResCode)
}

func Test_should_handle_empty_quiz_results(t *testing.T) {
	t.Parallel()

	emptyResults := dto.QuizResults{
		UserID:  uuid.New().String(),
		Results: []dto.WordQuizResult{},
	}
	httpResCode := postNoContent("/api/v1/app/quiz/results", ToJson(&emptyResults))
	assert.Equal(t, http.StatusOK, httpResCode)
}

func Test_should_handle_mixed_status_quiz_results(t *testing.T) {
	t.Parallel()

	// Create test word
	word := GenerateWord()
	var insertedWord models.Word
	httpResCode := post("/api/v1/tech/words", ToJson(&word), &insertedWord)
	assert.Equal(t, http.StatusCreated, httpResCode)

	// Submit mixed results for the same word
	statuses := []dto.ResultStatus{dto.Success, dto.Error, dto.Unanswered, dto.Success, dto.Success}
	for _, status := range statuses {
		results := dto.QuizResults{
			UserID: uuid.New().String(),
			Results: []dto.WordQuizResult{
				{WordID: insertedWord.ID, Status: status},
			},
		}
		httpResCode = postNoContent("/api/v1/app/quiz/results", ToJson(&results))
		assert.Equal(t, http.StatusOK, httpResCode)
		time.Sleep(100 * time.Millisecond) // Ensure different timestamps
	}
}
