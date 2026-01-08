package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	ModelName string `json:"model_name"`
	FilePath  string `json:"file_path"`
	Accuracy  string   `json:"accuracy"`
	Type      string `json:"type"`
	Active    bool   `json:active`
}

func RegisterModelRouter(r *gin.Engine) {
	r.GET("/api/v1/model", GetAllModelHandle)
	r.GET("/api/v1/model-active", GetAllActiveModelHandle)
	r.POST("/api/v1/model", CreateModelHandler)
	r.PUT("/api/v1/model/:id", UpdateModelHandler)
	r.DELETE("/api/v1/model/:id", DeleteModelHandler)
}

func GetAllModelHandle(c *gin.Context) {
	var models []Model
	result := db.Find(&models)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, models)
}

func GetAllActiveModelHandle(c *gin.Context) {
	var models []Model
	result := db.Where("active = ?", 1).Find(&models)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, models)
}

func CreateModelHandler(c *gin.Context) {
	var model Model
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

func UpdateModelHandler(c *gin.Context) {
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

func DeleteModelHandler(c *gin.Context) {
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
