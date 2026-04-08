package ll2datasyncer

import (
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

const PageSize = 100

type LaunchSyncer struct {
	*BaseSyncer
	ll2Service *ll2.LL2Service
	core       *core.MainService
	offset     int
	total      int
}

func NewLaunchSyncer(rl util.RateLimiter, ll2Service *ll2.LL2Service, core *core.MainService) *LaunchSyncer {
	ls := &LaunchSyncer{
		ll2Service: ll2Service,
		core:       core,
	}
	ls.BaseSyncer = NewBaseSyncer(rl, ls.sync)
	return ls
}

func (ls *LaunchSyncer) sync() {
	resp, err := ls.ll2Service.GetLaunchesFromAPI(PageSize, ls.offset)
	if err != nil {
		logrus.Errorf("get launches from ll2 api failed: %v", err)
		return
	}

	ls.total = resp.Count

	err = ls.ll2Service.SaveLaunchesToDB(resp.Results)
	if err != nil {
		logrus.Errorf("save launches to db failed: %v", err)
		return
	}

	err = ls.core.GenerateLaunchFromLL2(GenIdsFromLaunches(resp.Results))
	if err != nil {
		logrus.Errorf("generate launches from ll2 ids failed: %v", err)
		return
	}

	ls.offset += len(resp.Results)

	if ls.offset >= ls.total {
		ls.requestStop()
	}
}

func GenIdsFromLaunches(launches []*models.LL2LaunchDetailed) []string {
	ids := make([]string, 0, len(launches))
	for _, launch := range launches {
		ids = append(ids, launch.ID)
	}
	return ids
}
