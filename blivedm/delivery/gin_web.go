package delivery

import (
	webapi "github.com/AynaLivePlayer/blivedm-go/api"
	sgin "github.com/rhine-tech/scene/scenes/gin"
)

type dmInfoResponse struct {
	UID       int               `json:"uid"`
	DanmuInfo *webapi.DanmuInfo `json:"danmu_info"`
	Error     string            `json:"error"`
}

type dmInfoRequest struct {
	sgin.RequestQuery
	RoomId int `form:"room_id" binding:"required,number" json:"room_id"`
}

func (d *dmInfoRequest) Process(ctx *sgin.Context[*ginApp]) (data any, err error) {
	uid, danmuInfo, err := ctx.App.dmSrv.GetDanmuInfo(d.RoomId)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return &dmInfoResponse{
		UID:       uid,
		DanmuInfo: danmuInfo,
		Error:     errMsg}, nil
}
