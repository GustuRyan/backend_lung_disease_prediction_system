package handler

import (
	"net/http"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/domain"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/httpx"
	"github.com/gin-gonic/gin"
)

type datasetRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	DatasetName string `json:"dataset_name" binding:"required,max=160"`
	FilePath    string `json:"file_path" binding:"required,max=1024"`
	Description string `json:"description" binding:"max=10000"`
}

func (h *Handler) listDatasets(c *gin.Context) {
	items, err := h.service.Datasets(c.Request.Context())
	if err != nil {
		httpx.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) createDataset(c *gin.Context) {
	var request datasetRequest
	if !bind(c, &request) {
		return
	}
	item := domain.Dataset{UserID: request.UserID, DatasetName: request.DatasetName, FilePath: request.FilePath, Description: request.Description}
	if err := h.service.Create(c.Request.Context(), &item); err != nil {
		httpx.PersistenceError(c, err, "dataset")
		return
	}
	created(c, "dataset", item)
}

func (h *Handler) updateDataset(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	var request datasetRequest
	if !bind(c, &request) {
		return
	}
	err := h.service.Update(c.Request.Context(), &domain.Dataset{}, id, map[string]any{
		"user_id":      request.UserID,
		"dataset_name": request.DatasetName,
		"file_path":    request.FilePath,
		"description":  request.Description,
	})
	if err != nil {
		httpx.PersistenceError(c, err, "dataset")
		return
	}
	updated(c, "dataset")
}

func (h *Handler) deleteDataset(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), &domain.Dataset{}, id); err != nil {
		httpx.PersistenceError(c, err, "dataset")
		return
	}
	deleted(c, "dataset", id)
}
