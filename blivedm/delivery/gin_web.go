package delivery

import (
	webapi "github.com/AynaLivePlayer/blivedm-go/api"
	"github.com/gin-gonic/gin"
	"github.com/rhine-tech/scene"
	authMdw "github.com/rhine-tech/scene/lens/authentication/delivery/middleware"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"net/http"
)

type dmInfoResponse struct {
	UID       int               `json:"uid"`
	DanmuInfo *webapi.DanmuInfo `json:"danmu_info"`
	Error     string            `json:"error"`
}

type dmInfoV1Request struct {
	sgin.BaseAction
	sgin.RequestQuery
	RoomId int `form:"room_id" binding:"required,number" json:"room_id"`
}

func (d *dmInfoV1Request) GetRoute() scene.HttpRouteInfo {
	return scene.HttpRouteInfo{
		Method: http.MethodGet,
		Path:   "/web/v1/dm_info",
	}
}

func (d *dmInfoV1Request) Process(ctx *sgin.Context[*appContext]) (data any, err error) {
	return ctx.App.dmSrv.GetDanmuInfoCompatible(d.RoomId)
}

type dmInfoRequest struct {
	sgin.BaseAction
	sgin.RequestQuery
	RoomId int `form:"room_id" binding:"required,number" json:"room_id"`
}

func (d *dmInfoRequest) GetRoute() scene.HttpRouteInfo {
	return scene.HttpRouteInfo{
		Method: http.MethodGet,
		Path:   "/web/dm_info",
	}
}

func (d *dmInfoRequest) Process(ctx *sgin.Context[*appContext]) (data any, err error) {
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

type connLogRequest struct {
	sgin.BaseAction
	sgin.RequestQuery
	Offset int `form:"offset,default=0" binding:"number" json:"offset"`
	Limit  int `form:"limit,default=20" binding:"number" json:"limit"`
}

func (d *connLogRequest) GetRoute() scene.HttpRouteInfo {
	return scene.HttpRouteInfo{
		Method: http.MethodGet,
		Path:   "/connlog",
	}
}

func (r *connLogRequest) Middleware() gin.HandlersChain {
	return []gin.HandlerFunc{authMdw.GinRequireAuth(nil)}
}

func (d *connLogRequest) Process(ctx *sgin.Context[*appContext]) (data any, err error) {
	return ctx.App.connlog.ListEntries(int64(d.Offset), int64(d.Limit))
}

type roomLogRequest struct {
	sgin.RequestQuery
	Offset int `form:"offset,default=0" binding:"number" json:"offset"`
	Limit  int `form:"limit,default=20" binding:"number" json:"limit"`
}

func (r *roomLogRequest) GetRoute() scene.HttpRouteInfo {
	return scene.HttpRouteInfo{
		Method: http.MethodGet,
		Path:   "/roomlog",
	}
}

func (r *roomLogRequest) Middleware() gin.HandlersChain {
	return []gin.HandlerFunc{authMdw.GinRequireAuth(nil)}
}

func (r *roomLogRequest) Process(ctx *sgin.Context[*appContext]) (data any, err error) {
	return ctx.App.connlog.GetRoomLog(int64(r.Offset), int64(r.Limit))
}
