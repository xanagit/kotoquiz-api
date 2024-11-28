package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/dto"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type WordControllerDto interface {
	ListWordsIDs(c *gin.Context)
	ReadDtoWord(c *gin.Context)
}

type WordDtoControllerImpl struct {
	WordDtoService services.WordDtoService
}

func (s *WordDtoControllerImpl) ListWordsIDs(c *gin.Context) {
	tagIds := getQueryParamList(c, "tags")
	levelNameIds := getQueryParamList(c, "levelNames")
	nb, _ := getQueryParamInt(c, "nb", 30)

	wordIdsList, err := s.WordDtoService.ListWordsIDs(tagIds, levelNameIds, nb)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wordIdsList)
}

func (s *WordDtoControllerImpl) ListDtoWords(c *gin.Context) {
	// TODO : Ajouter requête paginée
	ids := getQueryParamList(c, "ids") // Récupère les IDs depuis le paramètre de requête
	lang := getQueryParamLang(c)

	var words []*dto.WordDTO
	var err error

	if len(ids) > 0 {
		words, err = s.WordDtoService.ListWordsDtoByIDs(ids, lang)
	} else {
		words = []*dto.WordDTO{}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, words)
}

func (s *WordDtoControllerImpl) ReadDtoWord(c *gin.Context) {
	id := c.Param("id")
	lang := getQueryParamLang(c)

	wordDto, err := s.WordDtoService.ReadWord(id, lang)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wordDto)
}
