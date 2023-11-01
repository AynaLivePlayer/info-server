package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/errcode"
	"github.com/rhine-tech/scene/model"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/lyric"
	"net/http"
)

type ginApp struct {
	sgin.CommonApp
	srv lyric.LyricService `aperture:""`
}

func NewGinApp() sgin.GinApplication {
	return &ginApp{}
}

func (g *ginApp) Name() scene.AppName {
	return "lyric.delivery.gin"
}

func (g *ginApp) Prefix() string {
	return "lyric"
}

func (g *ginApp) Create(engine *gin.Engine, router gin.IRouter) error {
	router.GET("/search", g.search)
	return nil
}

type paramSearch struct {
	Title  string `form:"title" json:"title" binding:"required"`
	Artist string `form:"artist" json:"artist"`
}

func (g *ginApp) search(c *gin.Context) {
	var param paramSearch
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorCodeResponse(errcode.ParameterError.WithDetail(err)))
		return
	}
	result, err := g.srv.Search(param.Title, param.Artist)
	if err != nil {
		c.JSON(200, model.TryErrorCodeResponse(err))
		return
	}
	c.JSON(200, model.NewDataResponse(result))
}
