package ll2datasyncer

import (
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

type AgencySyncer struct {
	*BaseSyncer
	ll2Service *ll2.LL2Service
	core       *core.MainService
	offset     int
	total      int
}

func NewAgencySyncer(rl util.RateLimiter, ll2Service *ll2.LL2Service, core *core.MainService) *AgencySyncer {
	as := &AgencySyncer{
		ll2Service: ll2Service,
		core:       core,
	}
	as.BaseSyncer = NewBaseSyncer(rl, as.sync)
	return as
}

func (as *AgencySyncer) sync() {
	resp, err := as.ll2Service.GetAngecyFromAPI(PageSize, as.offset)
	if err != nil {
		logrus.Errorf("get agencies from ll2 api failed: %v", err)
		return
	}

	as.total = resp.Count

	err = as.ll2Service.SaveAgenciesToDB(resp.Results)
	if err != nil {
		logrus.Errorf("save agencies to db failed: %v", err)
		return
	}

	err = as.core.GenerateAgencyFromLL2(GenIdsFromAgencies(resp.Results))
	if err != nil {
		logrus.Errorf("generate agencies from ll2 ids failed: %v", err)
		return
	}

	as.offset += len(resp.Results)

	if as.offset >= as.total {
		as.requestStop()
	}
}

func GenIdsFromAgencies(agencies []*models.LL2AgencyDetailed) []int64 {
	ids := make([]int64, 0, len(agencies))
	for _, agency := range agencies {
		ids = append(ids, int64(agency.ID))
	}
	return ids
}
