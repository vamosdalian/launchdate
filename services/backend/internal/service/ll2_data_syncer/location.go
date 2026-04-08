package ll2datasyncer

import (
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

type LocationSyncer struct {
	*BaseSyncer
	ll2Service *ll2.LL2Service
	core       *core.MainService
	offset     int
	total      int
}

func NewLocationSyncer(rl util.RateLimiter, ll2Service *ll2.LL2Service, core *core.MainService) *LocationSyncer {
	ls := &LocationSyncer{
		ll2Service: ll2Service,
		core:       core,
	}
	ls.BaseSyncer = NewBaseSyncer(rl, ls.sync)
	return ls
}

func (ls *LocationSyncer) sync() {
	resp, err := ls.ll2Service.GetLocationsFromApi(PageSize, ls.offset)
	if err != nil {
		logrus.Errorf("get locations from ll2 api failed: %v", err)
		return
	}

	ls.total = resp.Count

	var toSave []*models.LL2LocationSerializerWithPads
	for i := range resp.Results {
		toSave = append(toSave, &resp.Results[i])
	}

	err = ls.ll2Service.SaveLocationsToDB(toSave)
	if err != nil {
		logrus.Errorf("save locations to db failed: %v", err)
		return
	}

	err = ls.core.GenerateLaunchBaseFromLL2(GenIdsFromLocations(toSave))
	if err != nil {
		logrus.Errorf("generate launch base from ll2 ids failed: %v", err)
		return
	}

	ls.offset += len(resp.Results)

	if ls.offset >= ls.total {
		ls.requestStop()
	}
}

func GenIdsFromLocations(locations []*models.LL2LocationSerializerWithPads) []int64 {
	ids := make([]int64, 0, len(locations))
	for _, location := range locations {
		ids = append(ids, int64(location.ID))
	}
	return ids
}
