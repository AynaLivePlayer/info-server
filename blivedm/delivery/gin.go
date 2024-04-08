package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/rhine-tech/scene"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/blivedm"
)

type ginApp struct {
	//dmSrv     blivedm.BiliLiveDanmuInfoService `aperture:""`
	openblive blivedm.OpenBLiveApiService `aperture:""`
}

func (g *ginApp) Name() scene.ImplName {
	return scene.NewAppImplNameNoVer(blivedm.ModuleName, "gin")
}

func (g *ginApp) Prefix() string {
	return blivedm.ModuleName
}

func (g *ginApp) Create(engine *gin.Engine, router gin.IRouter) error {

	R := sgin.RequestWrapper(g)

	//router.GET("/casual/casual", R(&dmInfoRequest{}))

	router.GET("/openblive/app_start", R(&openbliveAppStartRequest{}))
	router.GET("/openblive/app_end", R(&openbliveAppEndRequest{}))
	router.GET("/openblive/heartbeat", R(&openbliveHeartBeatRequest{}))

	return nil
}

func (g *ginApp) Destroy() error {
	return nil
}

func NewGinApp(
	dmSrv blivedm.BiliLiveDanmuInfoService,
	openblive blivedm.OpenBLiveApiService) sgin.GinApplication {
	app := &ginApp{
		//dmSrv:     dmSrv,
		openblive: openblive,
	}
	return app
}

//type dmInfoRequest struct {
//	sgin.RequestQuery
//	RoomId string `form:"room_id" binding:"required,number" json:"room_id"`
//}
//
//func (d *dmInfoRequest) Process(ctx *sgin.Context[*ginApp]) (data any, err error) {
//	return ctx.App.dmSrv.GetDanmuInfo(cast.ToInt(d.RoomId))
//}
