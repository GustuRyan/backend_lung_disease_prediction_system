package handler

import (
	"net/http"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/domain"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/httpx"
	"github.com/gin-gonic/gin"
)

type modelRequest struct {
	ModelName string  `json:"model_name" binding:"required,max=160"`
	FilePath  string  `json:"file_path" binding:"required,max=1024"`
	Accuracy  float64 `json:"accuracy" binding:"gte=0,lte=1"`
	Type      string  `json:"type" binding:"required,max=80"`
	Active    bool    `json:"active"`
}

func (h *Handler) listModels(c *gin.Context) {
	h.respondModels(c, false)
}

func (h *Handler) listActiveModels(c *gin.Context) {
	h.respondModels(c, true)
}

func (h *Handler) respondModels(c *gin.Context, activeOnly bool) {
	items, err := h.service.Models(c.Request.Context(), activeOnly)
	if err != nil {
		httpx.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) createModel(c *gin.Context) {
	var request modelRequest
	if !bind(c, &request) {
		return
	}
	item := domain.Model{ModelName: request.ModelName, FilePath: request.FilePath, Accuracy: request.Accuracy, Type: request.Type, Active: request.Active}
	if err := h.service.Create(c.Request.Context(), &item); err != nil {
		httpx.PersistenceError(c, err, "model")
		return
	}
	created(c, "model", item)
}

func (h *Handler) updateModel(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	var request modelRequest
	if !bind(c, &request) {
		return
	}
	err := h.service.Update(c.Request.Context(), &domain.Model{}, id, map[string]any{
		"model_name": request.ModelName,
		"file_path":  request.FilePath,
		"accuracy":   request.Accuracy,
		"type":       request.Type,
		"active":     request.Active,
	})
	if err != nil {
		httpx.PersistenceError(c, err, "model")
		return
	}
	updated(c, "model")
}

func (h *Handler) deleteModel(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), &domain.Model{}, id); err != nil {
		httpx.PersistenceError(c, err, "model")
		return
	}
	deleted(c, "model", id)
}
