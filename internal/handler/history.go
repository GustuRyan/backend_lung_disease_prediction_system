package handler

import (
	"net/http"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/domain"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/httpx"
	"github.com/gin-gonic/gin"
)

type historyRequest struct {
	UserID           uint    `json:"user_id" binding:"required"`
	DiseaseID        uint    `json:"disease_id" binding:"required"`
	ImagePath        string  `json:"image_path" binding:"required,max=1024"`
	PredictionResult string  `json:"prediction_result" binding:"required,max=160"`
	ConfidenceResult float64 `json:"confidence_result" binding:"gte=0,lte=1"`
}

func (h *Handler) listHistory(c *gin.Context) {
	userID, ok := pathID(c, "id")
	if !ok {
		return
	}
	items, err := h.service.Histories(c.Request.Context(), userID)
	if err != nil {
		httpx.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) createHistory(c *gin.Context) {
	var request historyRequest
	if !bind(c, &request) {
		return
	}
	item := domain.History{
		UserID: request.UserID, DiseaseID: request.DiseaseID, ImagePath: request.ImagePath,
		PredictionResult: request.PredictionResult, ConfidenceResult: request.ConfidenceResult,
	}
	if err := h.service.Create(c.Request.Context(), &item); err != nil {
		httpx.PersistenceError(c, err, "history")
		return
	}
	created(c, "history", item)
}

func (h *Handler) updateHistory(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	var request historyRequest
	if !bind(c, &request) {
		return
	}
	err := h.service.Update(c.Request.Context(), &domain.History{}, id, map[string]any{
		"user_id":           request.UserID,
		"disease_id":        request.DiseaseID,
		"image_path":        request.ImagePath,
		"prediction_result": request.PredictionResult,
		"confidence_result": request.ConfidenceResult,
	})
	if err != nil {
		httpx.PersistenceError(c, err, "history")
		return
	}
	updated(c, "history")
}

func (h *Handler) deleteHistory(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), &domain.History{}, id); err != nil {
		httpx.PersistenceError(c, err, "history")
		return
	}
	deleted(c, "history", id)
}
