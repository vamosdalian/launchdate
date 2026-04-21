package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/subscription"
)

func (h *Handler) Subscribe(c *gin.Context) {
	var req models.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("invalid subscribe request body")
		h.Error(c, "invalid request body")
		return
	}

	if err := h.subscription.Subscribe(c.Request.Context(), req.Email); err != nil {
		h.logger.WithError(err).WithField("email", req.Email).Error("subscription request failed")
		switch {
		case errors.Is(err, subscription.ErrInvalidEmail):
			h.Error(c, "invalid email address")
		case errors.Is(err, subscription.ErrEmailNotConfigured):
			h.Error(c, "subscription is temporarily unavailable")
		default:
			h.Error(c, "failed to subscribe email")
		}
		return
	}

	h.Json(c, gin.H{"status": models.SubscriptionStatusSubscribed})
}

func (h *Handler) Unsubscribe(c *gin.Context) {
	token := c.Query("token")
	status, err := h.subscription.UnsubscribeByToken(c.Request.Context(), token)
	if err != nil {
		h.logger.WithError(err).WithField("token_present", token != "").Error("unsubscribe request failed")
		invalidToken := errors.Is(err, subscription.ErrInvalidToken)
		page := h.subscription.RenderUnsubscribeResultPage(status, invalidToken)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
		return
	}

	page := h.subscription.RenderUnsubscribeResultPage(status, false)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
}
