package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xanagit/kotoquiz-api/services"
	"net/http"
)

type LabelController interface {
	ListLabels(c *gin.Context)
}

type LabelControllerImpl struct {
	Service services.LabelService
}

func (lc *LabelControllerImpl) ListTags(c *gin.Context) {
	labels, err := lc.Service.ListLabels("TAG")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, labels)
}

func (lc *LabelControllerImpl) ListCategories(c *gin.Context) {
	labels, err := lc.Service.ListLabels("CATEGORY")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, labels)
}

func (lc *LabelControllerImpl) ListLevelNames(c *gin.Context) {
	labels, err := lc.Service.ListLabels("LEVEL_NAME")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, labels)
}
