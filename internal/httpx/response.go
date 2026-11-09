package httpx

import (
	"errors"
	"net/http"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/repository"
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

func InternalError(c *gin.Context, err error) {
	_ = c.Error(err)
	Error(c, http.StatusInternalServerError, "internal_error", "an unexpected error occurred")
}

func PersistenceError(c *gin.Context, err error, resource string) {
	if errors.Is(err, repository.ErrNotFound) {
		Error(c, http.StatusNotFound, "not_found", resource+" not found")
		return
	}
	if errors.Is(err, repository.ErrConflict) {
		Error(c, http.StatusConflict, "conflict", resource+" already exists")
		return
	}
	InternalError(c, err)
}
