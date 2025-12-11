package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID;references:ID"`

	DiseaseID uint    `json:"user_id"`
	Disease   Disease `gorm:"foreignKey:UserID;references:ID"`

	ImagePath string `json:"image_path"`
	PredictionResult string `json:"prediction_result"`
	ConfidenceResult string `json:"confidence_result"`
}

func RegisterHistoryRouter(r *gin.Engine) {
	
}

func GetAllHistoryHandler(c *gin.Context) {
	var historys []History
	result := db.Preload("User").Preload("Disease").Find(&historys)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, historys)
}

func CreateHistoryHandler(c *gin.Context) {
	var history []History
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&history)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History created", "data": history})
}

func UpdateHistoryHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var history History
	result := db.First(&history, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&history)
	c.JSON(http.StatusOK, gin.H{"message": "History updated", "data": history})
}

func DeleteHistoryHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var history History
	result := db.First(&history, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	db.Delete(&history)
	c.JSON(http.StatusOK, gin.H{"message": "History deleted", "id": id})
}