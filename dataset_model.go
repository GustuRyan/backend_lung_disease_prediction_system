package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DatasetModel struct {
	gorm.Model
	DatasetID uint    `json:"dataset_id"`
	Dataset   Dataset `gorm:"foreignKey:DatasetID;references:ID"`
	ModelID   uint    `json:"model_id"`
	Models    Model   `gorm:"foreignKey:ModelID;references:ID"`
}

func RegisterDatasetModelRouter(r *gin.Engine) {
	r.GET("/api/v1/dataset-model", GetAllDatasetModelHandle)
	r.POST("/api/v1/dataset-model", CreateDatasetModelHandler)
	r.PUT("/api/v1/dataset-model/:id", UpdateDatasetModelHandler)
	r.DELETE("/api/v1/dataset-model/:id", DeleteDatasetModelHandler)
}

func GetAllDatasetModelHandle(c *gin.Context) {
	var models []DatasetModel
	result := db.Preload("Dataset").Preload("Model").Find(&models)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, models)
}

func CreateDatasetModelHandler(c *gin.Context) {
	var model DatasetModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&model)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model created", "data": model})
}

func UpdateDatasetModelHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var model Model
	result := db.First(&model, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dataset not found"})
		return
	}

	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&model)
	c.JSON(http.StatusOK, gin.H{"message": "Model updated", "data": model})
}

func DeleteDatasetModelHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var model Model
	result := db.First(&model, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dataset not found"})
		return
	}

	db.Delete(&model)
	c.JSON(http.StatusOK, gin.H{"message": "Model deleted", "id": id})
}
