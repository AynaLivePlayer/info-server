package delivery

import (
	openblive "github.com/aynakeya/open-bilibili-live"
	sgin "github.com/rhine-tech/scene/scenes/gin"
)

type openbliveResponse struct {
	Result *openblive.AppStartResult `json:"result"`
	Error  *openblive.PublicError    `json:"error"`
}

type openbliveAppStartRequest struct {
	sgin.RequestQuery
	Code  string `form:"code" binding:"required" json:"code"`
	AppId int64  `form:"app_id" binding:"required,number" json:"app_id"`
}

func (d *openbliveAppStartRequest) Process(ctx *sgin.Context[*ginApp]) (data any, err error) {
	r, e := ctx.App.openblive.AppStart(d.Code, d.AppId)
	return openbliveResponse{Result: r, Error: e}, nil
}

type openbliveAppEndRequest struct {
	sgin.RequestQuery
	AppId  int64  `form:"app_id" binding:"required,number" json:"app_id"`
	GameId string `form:"game_id" binding:"required" json:"game_id"`
}

func (d *openbliveAppEndRequest) Process(ctx *sgin.Context[*ginApp]) (data any, err error) {
	e := ctx.App.openblive.AppEnd(d.AppId, d.GameId)
	return openbliveResponse{Result: nil, Error: e}, nil
}

type openbliveHeartBeatRequest struct {
	sgin.RequestQuery
	GameId string `form:"game_id" binding:"required" json:"game_id"`
}

func (d *openbliveHeartBeatRequest) Process(ctx *sgin.Context[*ginApp]) (data any, err error) {
	e := ctx.App.openblive.HearBeat(d.GameId)
	return openbliveResponse{Result: nil, Error: e}, nil
}
