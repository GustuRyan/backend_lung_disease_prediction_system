package handler

import (
	"errors"
	"net/http"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/httpx"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/service"
	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=120"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) listUsers(c *gin.Context) {
	users, err := h.service.Users(c.Request.Context())
	if err != nil {
		httpx.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) register(c *gin.Context) {
	var request registerRequest
	if !bind(c, &request) {
		return
	}
	user, err := h.service.Register(c.Request.Context(), request.Name, request.Email, request.Password)
	if errors.Is(err, service.ErrEmailInUse) {
		httpx.Error(c, http.StatusConflict, "email_in_use", err.Error())
		return
	}
	if err != nil {
		httpx.InternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registration successful", "user": user})
}

func (h *Handler) login(c *gin.Context) {
	var request loginRequest
	if !bind(c, &request) {
		return
	}
	user, err := h.service.Login(c.Request.Context(), request.Email, request.Password)
	if errors.Is(err, service.ErrInvalidCredentials) {
		httpx.Error(c, http.StatusUnauthorized, "invalid_credentials", err.Error())
		return
	}
	if err != nil {
		httpx.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "user": user})
}
