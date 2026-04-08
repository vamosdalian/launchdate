package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	ll2datasyncer "github.com/vamosdalian/launchdate-backend/internal/service/ll2_data_syncer"
)

func (h *Handler) StartLL2LaunchUpdate(c *gin.Context) {
	h.startSyncTask(c, ll2datasyncer.SyncTypeLaunch, "LL2 launch update started")
}

func (h *Handler) StartLL2AngecyUpdate(c *gin.Context) {
	h.startSyncTask(c, ll2datasyncer.SyncTypeAgency, "LL2 agency update started")
}

func (h *Handler) StartLL2LauncherUpdate(c *gin.Context) {
	h.startSyncTask(c, ll2datasyncer.SyncTypeLauncher, "LL2 launcher update started")
}

func (h *Handler) StartLL2LauncherFamilyUpdate(c *gin.Context) {
	h.startSyncTask(c, ll2datasyncer.SyncTypeLauncherFamily, "LL2 launcher family update started")
}

func (h *Handler) StartLL2LocationUpdate(c *gin.Context) {
	h.startSyncTask(c, ll2datasyncer.SyncTypeLocation, "LL2 location update started")
}

func (h *Handler) StartLL2PadUpdate(c *gin.Context) {
	h.startSyncTask(c, ll2datasyncer.SyncTypePad, "LL2 pad update started")
}

func (h *Handler) startSyncTask(c *gin.Context, syncType, message string) {
	err := h.ll2syncer.InitSync(syncType)
	if err != nil {
		h.Error(c, err.Error())
		return
	}
	h.Success(c, message)
}

func (h *Handler) GetLL2Launches(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	launches, err := h.ll2Server.GetLaunchesFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get launches: "+err.Error())
		return
	}
	h.Json(c, launches)
}

func (h *Handler) GetLL2Angecy(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	angecies, err := h.ll2Server.GetAngecyFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get angecies: "+err.Error())
		return
	}
	h.Json(c, angecies)
}

func (h *Handler) GetLL2LauncherFamilies(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	families, err := h.ll2Server.GetLauncherFamiliesFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get launcher families: "+err.Error())
		return
	}
	h.Json(c, families)
}

func (h *Handler) GetLL2Launchers(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	launchers, err := h.ll2Server.GetLaunchersFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get launchers: "+err.Error())
		return
	}
	h.Json(c, launchers)
}

func (h *Handler) GetLL2Locations(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	locations, err := h.ll2Server.GetLocationsFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get locations: "+err.Error())
		return
	}
	h.Json(c, locations)
}

func (h *Handler) GetLL2Pads(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	pads, err := h.ll2Server.GetPadsFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get pads: "+err.Error())
		return
	}
	h.Json(c, pads)
}
