package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/errcode"
	"github.com/rhine-tech/scene/lens/infrastructure/logger"
	"github.com/rhine-tech/scene/model"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"github.com/spf13/cast"
	"infoserver/bilisrv"
	"net/http"
)

type ginApp struct {
	sgin.CommonApp
	dmSrv bilisrv.BiliLiveDanmuInfoService
}

func (g *ginApp) Name() scene.AppName {
	return scene.AppName("bilisrv.app.gin")
}

func (g *ginApp) Prefix() string {
	return "bilisrv"
}

func (g *ginApp) Create(engine *gin.Engine, router gin.IRouter) error {
	router.GET("/dminfo", g.handleGetDmInfo)
	g.AppStatus = scene.AppStatusRunning
	return nil
}

func NewGinApp(logger logger.ILogger,
	dmSrv bilisrv.BiliLiveDanmuInfoService) sgin.GinApplication {
	app := &ginApp{
		dmSrv: dmSrv,
	}
	app.Logger = logger.WithPrefix(string(app.Name()))
	return app
}

type dmInfoParam struct {
	RoomId string `form:"room_id" binding:"required,number" json:"room_id"`
}

func (g *ginApp) handleGetDmInfo(c *gin.Context) {
	var param dmInfoParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorCodeResponse(
			errcode.ParameterError.WithDetail(err)))
		return
	}
	info, err := g.dmSrv.GetDanmuInfo(cast.ToInt(param.RoomId))
	if err != nil {
		c.JSON(http.StatusOK, model.TryErrorCodeResponse(err))
		return
	}
	c.JSON(http.StatusOK, model.NewDataResponse(info))
}
