package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type Dataset struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID;references:ID"`

	DatasetName string `json:"dataset_name"`
	FilePath    string `json:"file_path"`
	Description string `json:"description"`
}

func RegisterDatasetRouter(r *gin.Engine) {
	r.GET("/api/v1/dataset", GetAllDatasetHandler)
	r.POST("/api/v1/dataset", CreateDatasetHandler)
	r.PUT("/api/v1/dataset/:id", UpdateDatasetHandler)
	r.DELETE("/api/v1/dataset/:id", DeleteDatasetHandler)
}

func GetAllDatasetHandler(c *gin.Context) {
	var datasets []Dataset
	result := db.Preload("User").Find(&datasets)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, datasets)
}

func CreateDatasetHandler(c *gin.Context) {
	var dataset Dataset
	if err := c.ShouldBindJSON(&dataset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&dataset)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dataset created", "data": dataset})
}

func UpdateDatasetHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var dataset Dataset
	result := db.First(&dataset, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	if err := c.ShouldBindJSON(&dataset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&dataset)
	c.JSON(http.StatusOK, gin.H{"message": "Dataset updated", "data": dataset})
}

func DeleteDatasetHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var dataset Dataset
	result := db.First(&dataset, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	db.Delete(&dataset)
	c.JSON(http.StatusOK, gin.H{"message": "Dataset deleted", "id": id})
}