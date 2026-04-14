package ll2datasyncer

import (
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

type LauncherFamilySyncer struct {
	*BaseSyncer
	ll2Service *ll2.LL2Service
	report     func(map[string]interface{})
	offset     int
	total      int
}

func NewLauncherFamilySyncer(rl util.RateLimiter, ll2Service *ll2.LL2Service, report func(map[string]interface{})) *LauncherFamilySyncer {
	lfs := &LauncherFamilySyncer{
		ll2Service: ll2Service,
		report:     report,
	}
	lfs.BaseSyncer = NewBaseSyncer(rl, lfs.sync)
	return lfs
}

func (lfs *LauncherFamilySyncer) sync() {
	resp, err := lfs.ll2Service.GetLauncherFamiliesFromAPI(PageSize, lfs.offset)
	if err != nil {
		logrus.Errorf("get launcher families from ll2 api failed: %v", err)
		return
	}

	lfs.total = resp.Count

	err = lfs.ll2Service.SaveLauncherFamiliesToDB(resp.Results)
	if err != nil {
		logrus.Errorf("save launcher families to db failed: %v", err)
		return
	}

	// No core generation needed for Launcher Families currently

	lfs.offset += len(resp.Results)
	lfs.notifyProgress()

	if lfs.offset >= lfs.total {
		lfs.requestStop()
	}
}

func (lfs *LauncherFamilySyncer) notifyProgress() {
	if lfs.report != nil {
		lfs.report(buildCountProgress(lfs.offset, lfs.total))
	}
}
