package ll2datasyncer

import (
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

type LauncherSyncer struct {
	*BaseSyncer
	ll2Service *ll2.LL2Service
	core       *core.MainService
	offset     int
	total      int
}

func NewLauncherSyncer(rl util.RateLimiter, ll2Service *ll2.LL2Service, core *core.MainService) *LauncherSyncer {
	ls := &LauncherSyncer{
		ll2Service: ll2Service,
		core:       core,
	}
	ls.BaseSyncer = NewBaseSyncer(rl, ls.sync)
	return ls
}

func (ls *LauncherSyncer) sync() {
	resp, err := ls.ll2Service.GetLaunchersFromAPI(PageSize, ls.offset)
	if err != nil {
		logrus.Errorf("get launchers from ll2 api failed: %v", err)
		return
	}

	ls.total = resp.Count

	var toSave []models.LL2LauncherConfigNormal
	for _, r := range resp.Results {
		toSave = append(toSave, r.LL2LauncherConfigNormal)
	}

	err = ls.ll2Service.SaveLaunchersToDB(toSave)
	if err != nil {
		logrus.Errorf("save launchers to db failed: %v", err)
		return
	}

	// Assuming Launchers correspond to Rockets in core service
	err = ls.core.GenerateRocketsFromLL2(GenIdsFromLaunchers(resp.Results))
	if err != nil {
		logrus.Errorf("generate rockets from ll2 ids failed: %v", err)
		return
	}

	ls.offset += len(resp.Results)

	if ls.offset >= ls.total {
		ls.requestStop()
	}
}

func GenIdsFromLaunchers(launchers []models.LL2LauncherConfigDetailed) []int64 {
	ids := make([]int64, 0, len(launchers))
	for _, launcher := range launchers {
		ids = append(ids, int64(launcher.ID))
	}
	return ids
}
