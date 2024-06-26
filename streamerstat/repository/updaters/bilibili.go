package updaters

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"infoserver/streamerstat"
	"time"
)

var r = resty.New()
var api = "https://api.live.bilibili.com/xlive/web-room/v1/index/getRoomBaseInfo?req_biz=video"

type Bilibili struct {
}

func (b *Bilibili) Platform() string {
	return "bilibili"
}

func (b *Bilibili) GetStatus(roomId string) (*streamerstat.StreamerStatus, error) {
	resp, err := r.R().
		SetQueryParam("room_ids", roomId).
		Get(api)
	if err != nil {
		return nil, err
	}
	result := gjson.ParseBytes(resp.Body())
	prefixPath := "data.by_room_ids." + roomId + "."
	status := &streamerstat.StreamerStatus{
		Platform:      b.Platform(),
		RoomId:        roomId,
		RoomTitle:     result.Get(prefixPath + "title").String(),
		Username:      result.Get(prefixPath + "uname").String(),
		UserID:        result.Get(prefixPath + "uid").String(),
		Followers:     int(result.Get(prefixPath + "attention").Int()),
		Category:      result.Get(prefixPath + "area_name").String(),
		IsStreaming:   result.Get(prefixPath + "live_status").Bool(),
		LastCheckTime: time.Now(),
		LiveUrl:       result.Get(prefixPath + "live_url").String(),
	}
	if status.RoomTitle == "" {
		return status, errors.New("fail to get room status: " + result.Get("message").String())
	}
	return status, nil
}
