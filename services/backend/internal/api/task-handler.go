package api

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type StartTaskRequest struct {
	Type string `json:"type" binding:"required"`
}

type TaskActionRequest struct {
	Action string `json:"action" binding:"required"`
}

func (h *Handler) StartTask(c *gin.Context) {
	var req StartTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, "invalid request body")
		return
	}

	syncType := strings.ToLower(strings.TrimSpace(req.Type))
	if err := h.ll2syncer.InitSync(syncType); err != nil {
		h.Error(c, err.Error())
		return
	}

	h.Success(c, "task started successfully")
}

func (h *Handler) GetTask(c *gin.Context) {
	h.Json(c, h.ll2syncer.GetCurrentTask())
}

func (h *Handler) TaskAction(c *gin.Context) {
	var req TaskActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, "invalid request body")
		return
	}

	action := strings.ToLower(strings.TrimSpace(req.Action))

	var err error
	switch action {
	case "pause":
		err = h.ll2syncer.PauseSync()
	case "resume":
		err = h.ll2syncer.ResumeSync()
	case "cancel":
		err = h.ll2syncer.CancelSync()
	default:
		h.Error(c, "unknown action")
		return
	}

	if err != nil {
		h.Error(c, err.Error())
		return
	}

	h.Success(c, "task action applied successfully")
}
