package service

import (
	"errors"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/infrastructure/logger"
	"infoserver/streamerstat"
	"time"
)

type stremerServiceImpl struct {
	repo             streamerstat.StreamerStatusRepository
	updaters         map[string]streamerstat.StreamerStatusUpdater
	log              logger.ILogger `aperture:""`
	skipUpdatePeriod time.Duration  // if less than this period. skip update
}

func (s *stremerServiceImpl) Setup() error {
	return nil
}

func (s *stremerServiceImpl) SrvImplName() scene.ImplName {
	return streamerstat.Lens.ImplNameNoVer("IStreamerStatsService")
}

func StreamerStatsService(repo streamerstat.StreamerStatusRepository, updaters ...streamerstat.StreamerStatusUpdater) streamerstat.IStreamerStatsService {
	updaterMap := make(map[string]streamerstat.StreamerStatusUpdater)
	for _, updater := range updaters {
		updaterMap[updater.Platform()] = updater
	}
	return &stremerServiceImpl{
		repo:             repo,
		updaters:         updaterMap,
		skipUpdatePeriod: time.Minute * 10,
	}
}

func (s *stremerServiceImpl) UpdateStatus(platform string, roomId string) (*streamerstat.StreamerStatus, error) {
	updater, exists := s.updaters[platform]
	if !exists {
		s.log.ErrorW("failed when updating status, not able to find corresponding platform",
			"platform", platform)
		return nil, streamerstat.ErrFailToUpdateStatus.WithDetailStr("platform not found")
	}

	// if current status found.
	// and last check time is still less than duration. skip
	status, err := s.GetStatus(platform, roomId)
	if err == nil {
		if status.LastCheckTime.Add(s.skipUpdatePeriod).After(time.Now()) {
			s.log.Warnf("last check time is %s, skip current check", status.LastCheckTime)
			return status, nil
		}
	}

	status, err = updater.GetStatus(roomId)
	if err != nil {
		s.log.ErrorW("updater checked error",
			"platform", platform,
			"roomId", roomId,
			"error", err)
		return nil, streamerstat.ErrFailToUpdateStatus.WithDetail(err)
	}

	err = s.repo.UpdateStatus(status)
	if err != nil {
		s.log.ErrorW("repository update status error",
			"platform", platform,
			"roomId", roomId,
			"error", err)
		return nil, streamerstat.ErrFailToUpdateStatus.WithDetail(err)
	}

	return status, nil
}

func (s *stremerServiceImpl) GetStatus(platform string, roomId string) (*streamerstat.StreamerStatus, error) {
	status, err := s.repo.GetStatus(platform, roomId)
	if err != nil {
		s.log.WarnW("failed to get status",
			"error", err,
			"platform", platform,
			"room_id", roomId)
	}
	if !errors.Is(err, streamerstat.ErrStatusNotFound) && err != nil {
		err = streamerstat.ErrGetStatusUnknownError
	}
	return status, err
}

func (s *stremerServiceImpl) GetStatusBatch(platform string, roomIds []string) ([]*streamerstat.StreamerStatus, error) {
	statuses, err := s.repo.GetStatusBatch(platform, roomIds)
	if err != nil {
		s.log.ErrorW("failed to get status", "error", err)
		return nil, err
	}
	// Create a map to quickly lookup statuses by room ID
	statusMap := make(map[string]*streamerstat.StreamerStatus)
	for _, status := range statuses {
		statusMap[status.RoomId] = status
	}
	// Create the result slice in the same order as roomIds
	result := make([]*streamerstat.StreamerStatus, len(roomIds))
	for i, roomId := range roomIds {
		if status, found := statusMap[roomId]; found {
			result[i] = status
		} else {
			result[i] = nil
		}
	}
	return result, nil
}
