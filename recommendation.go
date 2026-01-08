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
	r.GET("/api/v1/recommendations", GetAllRecommendationHandler)
	r.GET("/api/v1/recommendations/:disease_id", GetRecommendationByDiseaseHandler)
	r.POST("/api/v1/recommendations",  CreateRecommendationHandler)
	r.PUT("/api/v1/recommendations/:id", UpdateRecommendationHandler)
	r.DELETE("/api/v1/recommendations/:id", DeleteRecommendationHandler)
}

func GetAllRecommendationHandler(c *gin.Context) {
	var Recommendations []Recommendation
	result := db.Preload("Disease").Find(&Recommendations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, Recommendations)
}

func GetRecommendationByDiseaseHandler(c *gin.Context) {
	idStr := c.Param("disease_id")
	id, _ := strconv.Atoi(idStr)

	var Recommendations []Recommendation
	result := db.Where("disease_id = ?", id).Preload("Disease").Find(&Recommendations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, Recommendations)
}

func CreateRecommendationHandler(c *gin.Context) {
	var recommendation Recommendation
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