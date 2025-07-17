package delivery

import (
	"github.com/rhine-tech/scene"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"net/http"
)

var lxSourceMapping = map[string]string{
	"kw": "kuwo",
	"kg": "kugou",
	"tx": "qq",
	"wy": "netease",
}

type lxUrlRequest struct {
	sgin.BaseAction
	sgin.RequestQuery
	Provider string `json:"provider" form:"provider" binding:"required"`
	Id       string `json:"id" form:"id" binding:"required"`
	Quality  string `json:"quality" form:"quality" binding:"required"`
}

func (l *lxUrlRequest) GetRoute() scene.HttpRouteInfo {
	return scene.HttpRouteInfo{
		Method: http.MethodGet,
		Path:   "/lxmusic/url",
	}
}

func (l *lxUrlRequest) Process(ctx *sgin.Context[*appContext]) (data any, err error) {
	_, ok := lxSourceMapping[l.Provider]
	if ok {
		l.Provider = lxSourceMapping[l.Provider]
	}
	return ctx.App.srv.GetMediaUrl(l.Provider, l.Id, l.Quality)
}
