package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/rhine-tech/scene"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/lyric"
)

type ginApp struct {
	srv lyric.LyricService `aperture:""`
}

func NewGinApp() sgin.GinApplication {
	return &ginApp{}
}

func (g *ginApp) Destroy() error {
	return nil
}

func (g *ginApp) Name() scene.ImplName {
	return lyric.Lens.ImplNameNoVer("GinApplication")
}

func (g *ginApp) Prefix() string {
	return "lyric"
}

func (g *ginApp) Create(engine *gin.Engine, router gin.IRouter) error {
	R := sgin.RequestWrapper(g)
	router.GET("/get", R(&getReq{}))
	router.GET("/search", R(&searchReq{}))
	return nil
}

type getReq struct {
	sgin.RequestQuery
	Title  string `form:"title" json:"title" binding:"required"`
	Artist string `form:"artist" json:"artist"`
}

func (s *getReq) Process(ctx *sgin.Context[*ginApp]) (data any, err error) {
	result, err := ctx.App.srv.GetLyric(s.Title, s.Artist)
	return result, err
}

type searchReq struct {
	sgin.RequestQuery
	Keyword string `form:"keyword" json:"keyword" binding:"required"`
}

func (s *searchReq) Process(ctx *sgin.Context[*ginApp]) (data any, err error) {
	return ctx.App.srv.Search(s.Keyword)
}
