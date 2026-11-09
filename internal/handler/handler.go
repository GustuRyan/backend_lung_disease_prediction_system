package handler

import (
	"net/http"
	"strconv"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/httpx"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health/live", h.liveness)
	router.GET("/api/hello", func(c *gin.Context) { c.String(http.StatusOK, "Hello from LungCare API!") })

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", h.listUsers)
		v1.POST("/users", h.register)
		v1.POST("/register", h.register)
		v1.POST("/login", h.login)

		v1.GET("/disease", h.listDiseases)
		v1.GET("/disease/:name", h.getDiseaseByName)
		v1.POST("/disease", h.createDisease)
		v1.PUT("/disease/:id", h.updateDisease)
		v1.DELETE("/disease/:id", h.deleteDisease)

		v1.GET("/dataset", h.listDatasets)
		v1.POST("/dataset", h.createDataset)
		v1.PUT("/dataset/:id", h.updateDataset)
		v1.DELETE("/dataset/:id", h.deleteDataset)

		v1.GET("/model", h.listModels)
		v1.GET("/model-active", h.listActiveModels)
		v1.POST("/model", h.createModel)
		v1.PUT("/model/:id", h.updateModel)
		v1.DELETE("/model/:id", h.deleteModel)

		v1.GET("/recommendations", h.listRecommendations)
		v1.GET("/recommendations/:disease_id", h.listRecommendationsByDisease)
		v1.POST("/recommendations", h.createRecommendation)
		v1.PUT("/recommendations/:id", h.updateRecommendation)
		v1.DELETE("/recommendations/:id", h.deleteRecommendation)

		v1.GET("/history/:id", h.listHistory)
		v1.POST("/history", h.createHistory)
		v1.PUT("/history/:id", h.updateHistory)
		v1.DELETE("/history/:id", h.deleteHistory)

		v1.GET("/dataset-model", h.listDatasetModels)
		v1.POST("/dataset-model", h.createDatasetModel)
		v1.PUT("/dataset-model/:id", h.updateDatasetModel)
		v1.DELETE("/dataset-model/:id", h.deleteDatasetModel)
	}
}

func (h *Handler) liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func bind(c *gin.Context, destination any) bool {
	if err := c.ShouldBindJSON(destination); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid_request", err.Error())
		return false
	}
	return true
}

func pathID(c *gin.Context, parameter string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(parameter), 10, 64)
	if err != nil || value == 0 {
		httpx.Error(c, http.StatusBadRequest, "invalid_id", parameter+" must be a positive integer")
		return 0, false
	}
	return uint(value), true
}

func created(c *gin.Context, resource string, data any) {
	c.JSON(http.StatusCreated, gin.H{"message": resource + " created", "data": data})
}

func updated(c *gin.Context, resource string) {
	c.JSON(http.StatusOK, gin.H{"message": resource + " updated"})
}

func deleted(c *gin.Context, resource string, id uint) {
	c.JSON(http.StatusOK, gin.H{"message": resource + " deleted", "id": id})
}
