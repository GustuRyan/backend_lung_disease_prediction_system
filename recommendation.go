package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Recommendation struct {
	gorm.Model
	DiseaseID uint    `json:"disease_id"`
	Disease   Disease `gorm:"foreignKey:DiseaseID;references:ID"`

	Type               string `json:"type"`
	RecommendationText string `json:"recommendation_text"`
}

func RegisterRecommendationRouter(r *gin.Engine) {
	
}

func GetAllRecommendationHandler(c *gin.Context) {
	var Recommendations []Recommendation
	result := db.Preload("User").Find(&Recommendations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, Recommendations)
}

func CreateRecommendationHandler(c *gin.Context) {
	var recommendation []Recommendation
	if err := c.ShouldBindJSON(&recommendation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&recommendation)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recommendation created", "data": recommendation})
}

func UpdateRecommendationHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var recommendation Recommendation
	result := db.First(&recommendation, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recommendation not found"})
		return
	}

	if err := c.ShouldBindJSON(&recommendation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&recommendation)
	c.JSON(http.StatusOK, gin.H{"message": "Recommendation updated", "data": recommendation})
}

func DeleteRecommendationHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var recommendation Recommendation
	result := db.First(&recommendation, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recommendation not found"})
		return
	}

	db.Delete(&recommendation)
	c.JSON(http.StatusOK, gin.H{"message": "Recommendation deleted", "id": id})
}