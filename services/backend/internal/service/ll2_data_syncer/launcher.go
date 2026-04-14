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
	report     func(map[string]interface{})
	offset     int
	total      int
}

func NewLauncherSyncer(rl util.RateLimiter, ll2Service *ll2.LL2Service, core *core.MainService, report func(map[string]interface{})) *LauncherSyncer {
	ls := &LauncherSyncer{
		ll2Service: ll2Service,
		core:       core,
		report:     report,
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

	err = ls.ll2Service.SaveLaunchersToDB(resp.Results)
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
	ls.notifyProgress()

	if ls.offset >= ls.total {
		ls.requestStop()
	}
}

func (ls *LauncherSyncer) notifyProgress() {
	if ls.report != nil {
		ls.report(buildCountProgress(ls.offset, ls.total))
	}
}

func GenIdsFromLaunchers(launchers []models.LL2LauncherConfigDetailed) []int64 {
	ids := make([]int64, 0, len(launchers))
	for _, launcher := range launchers {
		ids = append(ids, int64(launcher.ID))
	}
	return ids
}
