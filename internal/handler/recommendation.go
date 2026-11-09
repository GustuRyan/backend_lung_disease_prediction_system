package handler

import (
	"net/http"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/domain"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/httpx"
	"github.com/gin-gonic/gin"
)

type recommendationRequest struct {
	DiseaseID          uint   `json:"disease_id" binding:"required"`
	Type               string `json:"type" binding:"required,max=80"`
	RecommendationText string `json:"recommendation_text" binding:"required,max=10000"`
}

func (h *Handler) listRecommendations(c *gin.Context) {
	h.respondRecommendations(c, nil)
}

func (h *Handler) listRecommendationsByDisease(c *gin.Context) {
	id, ok := pathID(c, "disease_id")
	if !ok {
		return
	}
	h.respondRecommendations(c, &id)
}

func (h *Handler) respondRecommendations(c *gin.Context, diseaseID *uint) {
	items, err := h.service.Recommendations(c.Request.Context(), diseaseID)
	if err != nil {
		httpx.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) createRecommendation(c *gin.Context) {
	var request recommendationRequest
	if !bind(c, &request) {
		return
	}
	item := domain.Recommendation{DiseaseID: request.DiseaseID, Type: request.Type, RecommendationText: request.RecommendationText}
	if err := h.service.Create(c.Request.Context(), &item); err != nil {
		httpx.PersistenceError(c, err, "recommendation")
		return
	}
	created(c, "recommendation", item)
}

func (h *Handler) updateRecommendation(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	var request recommendationRequest
	if !bind(c, &request) {
		return
	}
	err := h.service.Update(c.Request.Context(), &domain.Recommendation{}, id, map[string]any{
		"disease_id":          request.DiseaseID,
		"type":                request.Type,
		"recommendation_text": request.RecommendationText,
	})
	if err != nil {
		httpx.PersistenceError(c, err, "recommendation")
		return
	}
	updated(c, "recommendation")
}

func (h *Handler) deleteRecommendation(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), &domain.Recommendation{}, id); err != nil {
		httpx.PersistenceError(c, err, "recommendation")
		return
	}
	deleted(c, "recommendation", id)
}
