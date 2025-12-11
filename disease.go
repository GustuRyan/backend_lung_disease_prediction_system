package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Disease struct {
	gorm.Model
	DiseaseName string `json:"disease_name"`
	Description string `json:"description"`
}

func RegisterDiseaseRouter(r *gin.Engine) {
	
}

func GetAllDiseaseHandler(c *gin.Context) {
	var diseases []Disease
	if err := db.Find(&diseases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, diseases)
}

func CreateDiseaseHandler(c *gin.Context) {
	var disease []Disease
	if err := c.ShouldBindJSON(&disease); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&disease)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Disease created", "data": disease})
}

func UpdateDiseaseHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var disease Disease
	result := db.First(&disease, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Disease not found"})
		return
	}

	if err := c.ShouldBindJSON(&disease); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&disease)
	c.JSON(http.StatusOK, gin.H{"message": "Disease updated", "data": disease})
}

func DeleteDiseaseHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var disease Disease
	result := db.First(&disease, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Disease not found"})
		return
	}

	db.Delete(&disease)
	c.JSON(http.StatusOK, gin.H{"message": "Disease deleted", "id": id})
}