package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.Error(c, "failed to get file")
		return
	}
	defer file.Close()

	key, err := h.imageService.UploadImage(c.Request.Context(), file, header.Filename)
	if err != nil {
		h.logger.Errorf("failed to upload image: %v", err)
		h.Error(c, "failed to upload image")
		return
	}

	h.Json(c, gin.H{"key": key})
}

func (h *Handler) ListImages(c *gin.Context) {
	type ListQuery struct {
		Limit  int `form:"limit,default=10"`
		Offset int `form:"offset,default=0"`
	}
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.Error(c, "invalid query params")
		return
	}

	list, err := h.imageService.ListImages(query.Limit, query.Offset)
	if err != nil {
		h.logger.Errorf("failed to list images: %v", err)
		h.Error(c, "failed to list images")
		return
	}

	h.Json(c, list)
}

func (h *Handler) DeleteImage(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		h.Error(c, "key is required")
		return
	}

	if err := h.imageService.DeleteImage(key); err != nil {
		h.logger.Errorf("failed to delete image: %v", err)
		h.Error(c, "failed to delete image")
		return
	}

	h.Json(c, "ok")
}

func (h *Handler) GenerateThumb(c *gin.Context) {
	var req struct {
		ID     string `json:"id" binding:"required"`
		Width  int    `json:"width" binding:"required"`
		Height int    `json:"height" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, "invalid params")
		return
	}

	id, _ := strconv.ParseInt(req.ID, 10, 64)

	if err := h.imageService.GenerateThumb(c.Request.Context(), id, req.Width, req.Height); err != nil {
		h.logger.Errorf("failed to generate thumb: %v", err)
		h.Error(c, "failed to generate thumb")
		return
	}

	h.Json(c, "ok")
}
