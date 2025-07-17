package delivery

import (
	"github.com/rhine-tech/scene"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"net/http"
	"scene-service/meowsic"
)

type appContext struct {
	srv meowsic.IMeowsicService `aperture:""`
}

func GinApp() sgin.GinApplication {
	return &sgin.AppRoutes[appContext]{
		AppName:  meowsic.Lens.ImplNameNoVer("GinApplication"),
		BasePath: meowsic.Lens.String(),
		Actions: []sgin.Action[*appContext]{
			new(getStatusRequest),
			new(updateStatusRequest),
			new(getProvidersRequest),
			new(setProviderSessionRequest),
			new(removeProviderSessionRequest),
			new(lxUrlRequest),
		},
		Context: appContext{
			srv: nil,
		},
	}
}

type getStatusRequest struct {
	sgin.BaseAction
	sgin.RequestURI
	Provider string `uri:"provider" binding:"required"`
}

func (l *getStatusRequest) GetRoute() scene.HttpRouteInfo {
	return scene.HttpRouteInfo{
		Method: http.MethodGet,
		Path:   "/status/latest/:provider",
	}
}

func (l *getStatusRequest) Process(ctx *sgin.Context[*appContext]) (data any, err error) {
	return ctx.App.srv.WithContext(ctx).GetStatus(l.Provider)
}

type updateStatusRequest struct {
	sgin.BaseAction
	sgin.RequestURI
	Provider string `uri:"provider" binding:"required"`
}

func (l *updateStatusRequest) GetRoute() scene.HttpRouteInfo {
	return scene.HttpRouteInfo{
		Method: http.MethodGet,
		Path:   "/status/update/:provider",
	}
}

func (l *updateStatusRequest) Process(ctx *sgin.Context[*appContext]) (data any, err error) {
	return nil, ctx.App.srv.WithContext(ctx).UpdateStatus(l.Provider)
}

type getProvidersRequest struct {
	sgin.BaseAction
	sgin.RequestNoParam
}

func (l *getProvidersRequest) GetRoute() scene.HttpRouteInfo {
	return scene.HttpRouteInfo{
		Method: http.MethodGet,
		Path:   "/providers",
	}
}

func (l *getProvidersRequest) Process(ctx *sgin.Context[*appContext]) (data any, err error) {
	return ctx.App.srv.WithContext(ctx).ListAllProvider()
}

type setProviderSessionRequest struct {
	sgin.BaseAction
	sgin.RequestFormUrlEncoded
	Session string `form:"session" binding:"required"`
}

func (l *setProviderSessionRequest) GetRoute() scene.HttpRouteInfo {
	return scene.HttpRouteInfo{
		Method: http.MethodPut,
		Path:   "/login/session/:provider",
	}
}

func (l *setProviderSessionRequest) Process(ctx *sgin.Context[*appContext]) (data any, err error) {
	provider := ctx.Param("provider")
	return nil, ctx.App.srv.WithContext(ctx).RestoreLogin(provider, l.Session)
}

type removeProviderSessionRequest struct {
	sgin.BaseAction
	sgin.RequestNoParam
}

func (l *removeProviderSessionRequest) GetRoute() scene.HttpRouteInfo {
	return scene.HttpRouteInfo{
		Method: http.MethodDelete,
		Path:   "/login/session/:provider",
	}
}

func (l *removeProviderSessionRequest) Process(ctx *sgin.Context[*appContext]) (data any, err error) {
	provider := ctx.Param("provider")
	return nil, ctx.App.srv.WithContext(ctx).Logout(provider)
}
