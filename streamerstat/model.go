package streamerstat

import (
	"encoding/json"
	"time"
)

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

func (s StreamerStatus) MarshalJSON() ([]byte, error) {
	type Alias StreamerStatus
	return json.Marshal(&struct {
		LastCheckTime string `json:"last_check_time"`
		Alias
	}{
		LastCheckTime: s.LastCheckTime.Format("2006-01-02 15:04:05"),
		Alias:         (Alias)(s),
	})
}
