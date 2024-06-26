package streamerstat

import (
	"github.com/rhine-tech/scene"
	"time"
)

const Lens scene.ModuleName = "streamerstat"

type StreamerStatus struct {
	Platform      string    `json:"platform" gorm:"index"`
	RoomId        string    `json:"room_id" gorm:"index"`
	RoomTitle     string    `json:"room_title"`
	Username      string    `json:"username"`
	UserID        string    `json:"user_id"`
	Followers     int       `json:"followers"`
	Category      string    `json:"category"`
	IsStreaming   bool      `json:"is_streaming"`
	LastCheckTime time.Time `json:"last_check_time"`
	LiveUrl       string    `json:"live_url"`
}

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
	GetStatus(platform string, roomId string) (*StreamerStatus, error)
	// GetStatusBatch should return a list of *StreamerStatus
	//
	// the return array should have same size as roomId,
	// the order should also be the same as roomId parameter
	//
	// if corresponding status is not found, status should be nil
	GetStatusBatch(platform string, roomIds []string) ([]*StreamerStatus, error)
}
