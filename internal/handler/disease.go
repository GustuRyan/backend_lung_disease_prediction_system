package handler

import (
	"net/http"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/domain"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/httpx"
	"github.com/gin-gonic/gin"
)

type diseaseRequest struct {
	DiseaseName string `json:"disease_name" binding:"required,max=160"`
	Description string `json:"description" binding:"max=10000"`
}

func (h *Handler) listDiseases(c *gin.Context) {
	items, err := h.service.Diseases(c.Request.Context())
	if err != nil {
		httpx.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) getDiseaseByName(c *gin.Context) {
	item, err := h.service.DiseaseByName(c.Request.Context(), c.Param("name"))
	if err != nil {
		httpx.PersistenceError(c, err, "disease")
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) createDisease(c *gin.Context) {
	var request diseaseRequest
	if !bind(c, &request) {
		return
	}
	item := domain.Disease{DiseaseName: request.DiseaseName, Description: request.Description}
	if err := h.service.Create(c.Request.Context(), &item); err != nil {
		httpx.PersistenceError(c, err, "disease")
		return
	}
	created(c, "disease", item)
}

func (h *Handler) updateDisease(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	var request diseaseRequest
	if !bind(c, &request) {
		return
	}
	err := h.service.Update(c.Request.Context(), &domain.Disease{}, id, map[string]any{
		"disease_name": request.DiseaseName,
		"description":  request.Description,
	})
	if err != nil {
		httpx.PersistenceError(c, err, "disease")
		return
	}
	updated(c, "disease")
}

func (h *Handler) deleteDisease(c *gin.Context) {
	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), &domain.Disease{}, id); err != nil {
		httpx.PersistenceError(c, err, "disease")
		return
	}
	deleted(c, "disease", id)
}
