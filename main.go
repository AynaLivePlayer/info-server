package main

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/engines"
	config "github.com/rhine-tech/scene/lens/infrastructure/config/builder"
	"github.com/rhine-tech/scene/registry"
	sgin "github.com/rhine-tech/scene/scenes/gin"

	asynctask "github.com/rhine-tech/scene/lens/infrastructure/asynctask/builder"
	datasource "github.com/rhine-tech/scene/lens/infrastructure/datasource/builder"
	ingestion "github.com/rhine-tech/scene/lens/infrastructure/ingestion/builder"
	logger "github.com/rhine-tech/scene/lens/infrastructure/logger/builder"

	bilisrv "infoserver/bilisrv/builder"
	lyric "infoserver/lyric/builder"
)

func main() {
	config.Init(".env.info_server")
	builders := scene.BuilderArray{
		logger.Builder{},
		asynctask.Thunnus{},
		datasource.MysqlBuilderFromConfig("scene.mysql"),
		ingestion.DummyBuilder{},
		bilisrv.JsonDB{},
		lyric.DefaultBuilder{},
	}
	scene.BuildInitArray(builders).Inits()
	registry.Logger.Infof("starting info_server...")
	engine := engines.NewEngine(registry.Logger,
		sgin.NewAppContainer(
			registry.Config.GetString("scene.app.gin.addr"),
			scene.BuildApps[sgin.GinApplication](builders)...,
		))
	if err := engine.Run(); err != nil {
		registry.Logger.Errorf("engine error: %s", err)
	}
}
