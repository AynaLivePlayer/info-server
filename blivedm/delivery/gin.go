package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/rhine-tech/scene"
	authMdw "github.com/rhine-tech/scene/lens/authentication/delivery/middleware"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/blivedm"
)

type ginApp struct {
	dmSrv     blivedm.WebDanmuService      `aperture:""`
	openblive blivedm.OpenBLiveApiService  `aperture:""`
	connlog   blivedm.ConnectionLogService `aperture:""`
}

func (g *ginApp) Name() scene.ImplName {
	return blivedm.Lens.ImplNameNoVer("GinApplication")
}

func (g *ginApp) Prefix() string {
	return blivedm.Lens.String()
}

func (g *ginApp) Create(engine *gin.Engine, router gin.IRouter) error {

	R := sgin.RequestWrapper(g)

	router.GET("/web/dm_info", R(&dmInfoRequest{}))

	router.GET("/openblive/app_start", R(&openbliveAppStartRequest{}))
	router.GET("/openblive/app_end", R(&openbliveAppEndRequest{}))
	router.GET("/openblive/heartbeat", R(&openbliveHeartBeatRequest{}))

	router.GET("/connlog", authMdw.GinRequireAuth(nil), R(new(connLogRequest)))
	router.GET("/roomlog", authMdw.GinRequireAuth(nil), R(new(roomLogRequest)))

	return nil
}

func (g *ginApp) Destroy() error {
	return nil
}

func NewGinApp(
	dmSrv blivedm.WebDanmuService,
	openblive blivedm.OpenBLiveApiService,
	connlog blivedm.ConnectionLogService) sgin.GinApplication {
	app := &ginApp{
		dmSrv:     dmSrv,
		openblive: openblive,
		connlog:   connlog,
	}
	return app
}
