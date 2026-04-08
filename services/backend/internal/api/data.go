package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
)

func (h *Handler) GetLaunches(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortField := c.DefaultQuery("sort_by", "time")
	sortOrderParam := strings.ToLower(c.DefaultQuery("sort_order", "asc"))
	sortOrder := 1
	if sortOrderParam == "desc" {
		sortOrder = -1
	}
	query := core.LaunchQuery{
		Limit:         limit,
		Offset:        offset,
		Name:          c.Query("name"),
		Status:        c.Query("status"),
		LaunchService: c.Query("launch_service_provider"),
		Rocket:        c.Query("rocket"),
		Mission:       c.Query("mission"),
		SortBy:        sortField,
		SortOrder:     sortOrder,
	}
	launches, err := h.core.GetLaunches(query)
	if err != nil {
		h.Error(c, "failed to get launches: "+err.Error())
		return
	}
	h.Json(c, launches)
}

func (h *Handler) GetAgencies(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortField := c.DefaultQuery("sort_by", "")
	sortOrderParam := strings.ToLower(c.DefaultQuery("sort_order", "asc"))
	sortOrder := 1
	if sortOrderParam == "desc" {
		sortOrder = -1
	}
	query := core.AgencyQuery{
		Limit:     limit,
		Offset:    offset,
		Name:      c.Query("name"),
		Type:      c.Query("type"),
		Country:   c.Query("country"),
		SortBy:    sortField,
		SortOrder: sortOrder,
	}
	agencies, err := h.core.GetAgencies(query)
	if err != nil {
		h.Error(c, "failed to get agencies: "+err.Error())
		return
	}
	h.Json(c, agencies)
}

func (h *Handler) GetAgencyByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Error(c, "invalid agency id")
		return
	}

	agency, err := h.core.GetAgency(id)
	if err != nil {
		h.Error(c, "failed to get agency: "+err.Error())
		return
	}

	h.Json(c, agency)
}

func (h *Handler) UpdateAgency(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Error(c, "invalid agency id")
		return
	}

	var req models.Agency
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, "invalid request body")
		return
	}
	req.ID = id

	if err := h.core.UpdateAgency(&req); err != nil {
		h.Error(c, "failed to update agency: "+err.Error())
		return
	}

	h.Success(c, "agency updated successfully")
}

func (h *Handler) GetRockets(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortField := c.DefaultQuery("sort_by", "")
	sortOrderParam := strings.ToLower(c.DefaultQuery("sort_order", "asc"))
	sortOrder := 1
	if sortOrderParam == "desc" {
		sortOrder = -1
	}
	query := core.RocketQuery{
		Limit:     limit,
		Offset:    offset,
		FullName:  c.Query("full_name"),
		Name:      c.Query("name"),
		Variant:   c.Query("variant"),
		SortBy:    sortField,
		SortOrder: sortOrder,
	}
	rockets, err := h.core.GetRockets(query)
	if err != nil {
		h.Error(c, "failed to get rockets: "+err.Error())
		return
	}
	h.Json(c, rockets)
}

func (h *Handler) GetRocketByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Error(c, "invalid rocket id")
		return
	}

	rocket, err := h.core.GetRocket(id)
	if err != nil {
		h.Error(c, "failed to get rocket: "+err.Error())
		return
	}

	h.Json(c, rocket)
}

func (h *Handler) UpdateRocket(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Error(c, "invalid rocket id")
		return
	}

	var req models.Rocket
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, "invalid request body")
		return
	}
	req.ID = id

	if err := h.core.UpdateRocket(&req); err != nil {
		h.Error(c, "failed to update rocket: "+err.Error())
		return
	}

	h.Success(c, "rocket updated successfully")
}

func (h *Handler) GetLaunchBases(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortField := c.DefaultQuery("sort_by", "")
	sortOrderParam := strings.ToLower(c.DefaultQuery("sort_order", "asc"))
	sortOrder := 1
	if sortOrderParam == "desc" {
		sortOrder = -1
	}
	query := core.LaunchBaseQuery{
		Limit:         limit,
		Offset:        offset,
		Name:          c.Query("name"),
		CelestialBody: c.Query("celestial_body"),
		Country:       c.Query("country"),
		SortBy:        sortField,
		SortOrder:     sortOrder,
	}
	launchBases, err := h.core.GetLaunchBases(query)
	if err != nil {
		h.Error(c, "failed to get launch bases: "+err.Error())
		return
	}
	h.Json(c, launchBases)
}

func (h *Handler) GetLaunchByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Error(c, "invalid launch id")
		return
	}

	launch, err := h.core.GetLaunch(id)
	if err != nil {
		h.Error(c, "failed to get launch: "+err.Error())
		return
	}

	h.Json(c, launch)
}

func (h *Handler) UpdateLaunch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Error(c, "invalid launch id")
		return
	}

	var req models.Launch
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, "invalid request body")
		return
	}
	req.ID = id

	if err := h.core.UpdateLaunch(&req); err != nil {
		h.Error(c, "failed to update launch: "+err.Error())
		return
	}

	h.Success(c, "launch updated successfully")
}

func (h *Handler) GetLaunchBaseByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Error(c, "invalid launch base id")
		return
	}

	launchBase, err := h.core.GetLaunchBase(id)
	if err != nil {
		h.Error(c, "failed to get launch base: "+err.Error())
		return
	}

	h.Json(c, launchBase)
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.core.GetStats()
	if err != nil {
		h.Error(c, "failed to get stats: "+err.Error())
		return
	}

	ll2Stats, err := h.ll2Server.GetStats()
	if err != nil {
		h.Error(c, "failed to get ll2 stats: "+err.Error())
		return
	}
	stats.LL2 = ll2Stats

	h.Json(c, stats)
}

func (h *Handler) GetPublicLaunchByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Error(c, "invalid launch id")
		return
	}

	launch, err := h.core.GetPublicLaunch(id)
	if err != nil {
		h.Error(c, "failed to get launch: "+err.Error())
		return
	}

	h.Json(c, launch)
}

func (h *Handler) GetPublicLaunches(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	launches, err := h.core.GetPublicLaunches(page)
	if err != nil {
		h.Error(c, "failed to get launches: "+err.Error())
		return
	}
	h.Json(c, launches)
}

func (h *Handler) GetPublicRockets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	rockets, err := h.core.GetPublicRockets(page)
	if err != nil {
		h.Error(c, "failed to get rockets: "+err.Error())
		return
	}
	h.Json(c, rockets)
}

func (h *Handler) GetPublicRocketByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Error(c, "invalid rocket id")
		return
	}

	rocket, err := h.core.GetPublicRocket(id)
	if err != nil {
		h.Error(c, "failed to get rocket: "+err.Error())
		return
	}

	h.Json(c, rocket)
}
