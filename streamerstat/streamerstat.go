package streamerstat

import (
	"github.com/rhine-tech/scene"
)

const Lens scene.ModuleName = "streamerstat"

type StreamerStatusUpdater interface {
	Platform() string
	GetStatus(roomId string) (*StreamerStatus, error)
}

type StreamerStatusRepository interface {
	UpdateStatus(status *StreamerStatus) error
	GetStatus(platform string, roomId string) (*StreamerStatus, error)
	GetStatusBatch(platform string, roomIds []string) ([]*StreamerStatus, error)
}

type IStreamerStatsService interface {
	scene.Service
	UpdateStatus(platform string, roomId string) (*StreamerStatus, error)
	// GetStatus should return *StreamerStatus if found.
	// if not found. err will be ErrStatusNotFound
	GetStatus(platform string, roomId string) (*StreamerStatus, error)
	// GetStatusBatch should return a list of *StreamerStatus
	//
	// the return array should have same size as roomId,
	// the order should also be the same as roomId parameter
	//
	// if corresponding status is not found, status should be nil
	GetStatusBatch(platform string, roomIds []string) ([]*StreamerStatus, error)
}
