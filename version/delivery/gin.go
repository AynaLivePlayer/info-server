package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/model"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/version"
)

type ginApp struct {
	srv version.VersionService `aperture:""`
}

func NewGinApp(srv version.VersionService) sgin.GinApplication {
	return &ginApp{srv: srv}
}

func (g *ginApp) Destroy() error {
	return nil
}

func (g *ginApp) Name() scene.ImplName {
	return version.Lens.ImplNameNoVer("GinApplication")
}

func (g *ginApp) Prefix() string {
	return "version"
}

func (g *ginApp) Create(engine *gin.Engine, router gin.IRouter) error {
	R := sgin.RequestWrapper(g)
	router.GET("/check_update", R(new(checkUpdateReq)))
	router.GET("/latest", R(new(getLatestReq)))
	return nil
}

type checkUpdateReq struct {
	sgin.RequestQuery
	ClientVersion uint32 `form:"client_version" json:"client_version" binding:"required"`
}

func (s *checkUpdateReq) Process(ctx *sgin.Context[*ginApp]) (data interface{}, err error) {
	update, hasUpdate := ctx.App.srv.CheckUpdate(version.Version(s.ClientVersion))
	return model.JsonResponse{
		"has_update": hasUpdate,
		"latest":     update,
	}, nil
}

type getLatestReq struct {
	sgin.RequestNoParam
}

func (s *getLatestReq) Process(ctx *sgin.Context[*ginApp]) (data interface{}, err error) {
	return ctx.App.srv.GetLatest()
}
