package ll2datasyncer

import (
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

type PadSyncer struct {
	*BaseSyncer
	ll2Service *ll2.LL2Service
	offset     int
	total      int
}

func NewPadSyncer(rl util.RateLimiter, ll2Service *ll2.LL2Service) *PadSyncer {
	ps := &PadSyncer{
		ll2Service: ll2Service,
	}
	ps.BaseSyncer = NewBaseSyncer(rl, ps.sync)
	return ps
}

func (ps *PadSyncer) sync() {
	resp, err := ps.ll2Service.GetPadsFromAPI(PageSize, ps.offset)
	if err != nil {
		logrus.Errorf("get pads from ll2 api failed: %v", err)
		return
	}

	ps.total = resp.Count

	var toSave []*models.LL2Pad
	for i := range resp.Results {
		toSave = append(toSave, &resp.Results[i])
	}

	err = ps.ll2Service.SavePadsToDB(toSave)
	if err != nil {
		logrus.Errorf("save pads to db failed: %v", err)
		return
	}

	// No core generation needed for Pads currently

	ps.offset += len(resp.Results)

	if ps.offset >= ps.total {
		ps.requestStop()
	}
}
