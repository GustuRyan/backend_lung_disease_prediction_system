package handler

import (
	"net/http"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/domain"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/httpx"
	"github.com/gin-gonic/gin"
)

type datasetModelRequest struct {
	DatasetID uint `json:"dataset_id" binding:"required"`
	ModelID   uint `json:"model_id" binding:"required"`
}

func (h *Handler) listDatasetModels(c *gin.Context) {
	items, err := h.service.DatasetModels(c.Request.Context())
	if err != nil {
		httpx.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) createDatasetModel(c *gin.Context) {
	var request datasetModelRequest
	if !bind(c, &request) {
		return
	}
	item := domain.DatasetModel{DatasetID: request.DatasetID, ModelID: request.ModelID}
	if err := h.service.Create(c.Request.Context(), &item); err != nil {
		httpx.PersistenceError(c, err, "dataset-model association")
		return
	}
	created(c, "dataset-model association", item)
}

func (h *Handler) updateDatasetModel(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	var request datasetModelRequest
	if !bind(c, &request) {
		return
	}
	err := h.service.Update(c.Request.Context(), &domain.DatasetModel{}, id, map[string]any{
		"dataset_id": request.DatasetID,
		"model_id":   request.ModelID,
	})
	if err != nil {
		httpx.PersistenceError(c, err, "dataset-model association")
		return
	}
	updated(c, "dataset-model association")
}

func (h *Handler) deleteDatasetModel(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), &domain.DatasetModel{}, id); err != nil {
		httpx.PersistenceError(c, err, "dataset-model association")
		return
	}
	deleted(c, "dataset-model association", id)
}
