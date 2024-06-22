package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/model"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/version"
	"time"

	authMdw "github.com/rhine-tech/scene/lens/authentication/delivery/middleware"
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
	router.GET("/list", R(new(listReq)))
	router.POST("/upsert", authMdw.GinRequireAuth(nil), R(new(upsertReq)))
	return nil
}

type upsertReq struct {
	sgin.RequestJson
	Version     string `json:"version" binding:"required"`
	Note        string `json:"note" binding:"required"`
	ReleaseTime int64  `json:"release_time,default=0"`
}

func (s *upsertReq) Process(ctx *sgin.Context[*ginApp]) (data interface{}, err error) {
	if s.ReleaseTime == 0 {
		s.ReleaseTime = time.Now().Unix()
	}
	return ctx.App.srv.UpsertVersion(version.VersionInfo{
		Version:     version.VersionFromString(s.Version),
		Note:        s.Note,
		ReleaseTime: s.ReleaseTime,
	}), nil
}

type listReq struct {
	sgin.RequestQuery
	Limit  int `form:"limit,default=20" json:"limit" binding:""`
	Offset int `form:"offset,limit=0" json:"offset" binding:""`
}

func (s *listReq) Process(ctx *sgin.Context[*ginApp]) (data interface{}, err error) {
	return ctx.App.srv.ListVersions(int64(s.Offset), int64(s.Limit))
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
