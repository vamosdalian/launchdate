package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	ll2datasyncer "github.com/vamosdalian/launchdate-backend/internal/service/ll2_data_syncer"
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
	task, err := h.ll2syncer.InitSync(syncType)
	if err != nil {
		h.Error(c, err.Error())
		return
	}

	h.Json(c, task)
}

func (h *Handler) GetTask(c *gin.Context) {
	h.Json(c, h.ll2syncer.GetCurrentTask())
}

func (h *Handler) GetTaskHistory(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.Error(c, "invalid limit")
		return
	}

	tasks, err := h.ll2syncer.GetTaskHistory(limit)
	if err != nil {
		h.Error(c, "failed to get task history: "+err.Error())
		return
	}

	h.Json(c, tasks)
}

func (h *Handler) TaskAction(c *gin.Context) {
	var req TaskActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, "invalid request body")
		return
	}

	action := strings.ToLower(strings.TrimSpace(req.Action))

	var (
		task *ll2datasyncer.TaskInfo
		err  error
	)
	switch action {
	case "pause":
		task, err = h.ll2syncer.PauseSync()
	case "resume":
		task, err = h.ll2syncer.ResumeSync()
	case "cancel":
		task, err = h.ll2syncer.CancelSync()
	default:
		h.Error(c, "unknown action")
		return
	}

	if err != nil {
		h.Error(c, err.Error())
		return
	}

	h.Json(c, task)
}
