package ll2datasyncer

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

type UpcomingSyncer struct {
	*BaseSyncer
	ll2Service *ll2.LL2Service
	core       *core.MainService

	pendingAgencyIDs   []int
	pendingLauncherIDs []int
	pendingPadIDs      []int
	pendingLocationIDs []int

	launchLimiter util.RateLimiter
}

func NewUpcomingSyncer(rl util.RateLimiter, ll2Service *ll2.LL2Service, core *core.MainService) *UpcomingSyncer {
	s := &UpcomingSyncer{
		ll2Service:    ll2Service,
		core:          core,
		launchLimiter: util.NewRateLimit(24 * time.Hour),
	}
	s.BaseSyncer = NewBaseSyncer(rl, s.sync)
	return s
}

func (s *UpcomingSyncer) sync() {
	// 1. Process pending queues first (Priority)
	if len(s.pendingAgencyIDs) > 0 {
		id := s.pendingAgencyIDs[0]
		s.pendingAgencyIDs = s.pendingAgencyIDs[1:]
		s.syncAgency(id)
		return
	}
	if len(s.pendingLauncherIDs) > 0 {
		id := s.pendingLauncherIDs[0]
		s.pendingLauncherIDs = s.pendingLauncherIDs[1:]
		s.syncLauncher(id)
		return
	}
	if len(s.pendingPadIDs) > 0 {
		id := s.pendingPadIDs[0]
		s.pendingPadIDs = s.pendingPadIDs[1:]
		s.syncPad(id)
		return
	}
	if len(s.pendingLocationIDs) > 0 {
		id := s.pendingLocationIDs[0]
		s.pendingLocationIDs = s.pendingLocationIDs[1:]
		s.syncLocation(id)
		return
	}

	// 2. Fetch upcoming launches
	// Frequency check: once every 24 hours using rate limiter
	if !s.launchLimiter.Allow() {
		return
	}

	// Find latest launch in DB (Requirement 1)
	latestLaunch, err := s.ll2Service.GetLatestLaunchFromDB()
	var latestID string
	if err != nil {
		logrus.Infof("No latest launch in DB or error: %v", err)
	} else {
		// Just logging as per requirement.
		// "启动后寻找数据库最新的一条launch" -> Done.
		latestID = latestLaunch.ID
		logrus.WithFields(logrus.Fields{
			"launch_id":   latestLaunch.ID,
			"launch_name": latestLaunch.Name,
		}).Info("Latest launch in DB")
	}

	// Fetch next 100 (Requirement 2)
	// Passed latestID as per requirement to GetLaunchUpComingFromAPI
	launchesResp, err := s.ll2Service.GetLaunchUpComingFromAPI(latestID, 100, 0)
	if err != nil {
		logrus.Errorf("failed to get upcoming launches: %v", err)
		return
	}

	if launchesResp == nil || len(launchesResp.Results) == 0 {
		logrus.Info("No upcoming launches found")
		return
	}

	s.processLaunches(launchesResp.Results)
}

func (s *UpcomingSyncer) processLaunches(launches []*models.LL2LaunchDetailed) {
	for _, launch := range launches {
		// Requirement 3: Check Launch ID
		existingLaunch, _ := s.ll2Service.GetLaunchByIDFromDB(launch.ID)
		if existingLaunch != nil {
			continue
		}

		logrus.Infof("New upcoming launch found: %s", launch.Name)

		// Save launch
		if err := s.ll2Service.SaveLaunchesToDB([]*models.LL2LaunchDetailed{launch}); err != nil {
			logrus.Errorf("failed to save launch %s: %v", launch.ID, err)
		}

		// Generate core launch
		if err := s.core.GenerateLaunchFromLL2([]string{launch.ID}); err != nil {
			logrus.Errorf("failed to generate core launch %s: %v", launch.ID, err)
		}

		// Requirement 4: Check Agency
		agencyID := launch.LaunchServiceProvider.ID
		if agencyID != 0 {
			if exists, _ := s.ll2Service.GetAgencyByIDFromDB(agencyID); exists == nil {
				if !containsInt(s.pendingAgencyIDs, agencyID) {
					logrus.Infof("Missing agency %d for launch %s, adding to queue", agencyID, launch.Name)
					s.pendingAgencyIDs = append(s.pendingAgencyIDs, agencyID)
				}
			}
		}

		// Requirement 5: Check Rocket / Launcher Config
		launcherID := launch.Rocket.Configuration.ID
		if launcherID != 0 {
			if exists, _ := s.ll2Service.GetLauncherByIDFromDB(launcherID); exists == nil {
				if !containsInt(s.pendingLauncherIDs, launcherID) {
					logrus.Infof("Missing launcher %d for launch %s, adding to queue", launcherID, launch.Name)
					s.pendingLauncherIDs = append(s.pendingLauncherIDs, launcherID)
				}
			}
		}

		// Requirement 6: Check Pad and Location
		padID := launch.Pad.Id
		if padID != 0 {
			if exists, _ := s.ll2Service.GetPadByIDFromDB(padID); exists == nil {
				if !containsInt(s.pendingPadIDs, padID) {
					logrus.Infof("Missing pad %d for launch %s, adding to queue", padID, launch.Name)
					s.pendingPadIDs = append(s.pendingPadIDs, padID)
				}
			}
		}

		locationID := launch.Pad.Location.ID
		if locationID != 0 {
			if exists, _ := s.ll2Service.GetLocationByIDFromDB(locationID); exists == nil {
				if !containsInt(s.pendingLocationIDs, locationID) {
					logrus.Infof("Missing location %d for launch %s, adding to queue", locationID, launch.Name)
					s.pendingLocationIDs = append(s.pendingLocationIDs, locationID)
				}
			}
		}
	}
}

func (s *UpcomingSyncer) syncAgency(id int) {
	logrus.Infof("Syncing missing agency: %d", id)
	agency, err := s.ll2Service.GetAgencyFromAPI(id)
	if err != nil {
		logrus.Errorf("failed to get agency %d: %v", id, err)
		return
	}
	if err := s.ll2Service.SaveAgenciesToDB([]*models.LL2AgencyDetailed{agency}); err != nil {
		logrus.Errorf("failed to save agency %d: %v", id, err)
	}
	if err := s.core.GenerateAgencyFromLL2([]int64{int64(id)}); err != nil {
		logrus.Errorf("failed to generate core agency %d: %v", id, err)
	}
}

func (s *UpcomingSyncer) syncLauncher(id int) {
	logrus.Infof("Syncing missing launcher: %d", id)
	launcher, err := s.ll2Service.GetLauncherFromAPI(id)
	if err != nil {
		logrus.Errorf("failed to get launcher %d: %v", id, err)
		return
	}

	// Convert Detailed to Normal for saving
	toSave := []models.LL2LauncherConfigNormal{launcher.LL2LauncherConfigNormal}

	if err := s.ll2Service.SaveLaunchersToDB(toSave); err != nil {
		logrus.Errorf("failed to save launcher %d: %v", id, err)
	}
	if err := s.core.GenerateRocketsFromLL2([]int64{int64(id)}); err != nil {
		logrus.Errorf("failed to generate core rocket %d: %v", id, err)
	}
}

func (s *UpcomingSyncer) syncPad(id int) {
	logrus.Infof("Syncing missing pad: %d", id)
	pad, err := s.ll2Service.GetPadFromAPI(id)
	if err != nil {
		logrus.Errorf("failed to get pad %d: %v", id, err)
		return
	}
	if err := s.ll2Service.SavePadsToDB([]*models.LL2Pad{pad}); err != nil {
		logrus.Errorf("failed to save pad %d: %v", id, err)
	}
}

func (s *UpcomingSyncer) syncLocation(id int) {
	logrus.Infof("Syncing missing location: %d", id)
	loc, err := s.ll2Service.GetLocationFromAPI(id)
	if err != nil {
		logrus.Errorf("failed to get location %d: %v", id, err)
		return
	}
	if err := s.ll2Service.SaveLocationsToDB([]*models.LL2LocationSerializerWithPads{loc}); err != nil {
		logrus.Errorf("failed to save location %d: %v", id, err)
	}
	if err := s.core.GenerateLaunchBaseFromLL2([]int64{int64(id)}); err != nil {
		logrus.Errorf("failed to generate core launch base %d: %v", id, err)
	}
}

func containsInt(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
