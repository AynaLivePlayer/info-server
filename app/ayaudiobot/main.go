package main

import (
	"github.com/rhine-tech/scene"
	orm "github.com/rhine-tech/scene/composition/orm/factory"
	"github.com/rhine-tech/scene/engines"
	asynctask "github.com/rhine-tech/scene/infrastructure/asynctask/factory"
	config "github.com/rhine-tech/scene/infrastructure/config/factory"
	datasouce "github.com/rhine-tech/scene/infrastructure/datasource/factory"
	logger "github.com/rhine-tech/scene/infrastructure/logger/factory"
	authentication "github.com/rhine-tech/scene/lens/authentication/factory"
	"github.com/rhine-tech/scene/registry"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	blivedm "infoserver/blivedm/factory"
	lyric "infoserver/lyric/factory"
)

var configFile = "conf.ini"

func init() {
	if scene.GetEnvironment() == scene.EnvDevelopment {
		configFile = "conf.dev.ini"
	} else {
		configFile = "conf.ini"
	}
}

func main() {
	config.InitINI(configFile)
	builders := scene.ModuleFactoryArray{
		logger.ZapFactory{}.Default(),
		datasouce.Mysql{}.Default(),
		asynctask.Ants{},
		orm.GormMysql{},
		authentication.GinAppMysql{}.Default(),
		blivedm.App{}.Default(),
		lyric.DefaultBuilder{},
	}
	scene.BuildInitArray(builders).Inits()
	registry.Logger.Infof("using config file: %s", configFile)
	engine := engines.NewEngine(registry.Logger,
		sgin.NewAppContainerWithPrefix(
			registry.Config.GetString("scene.app.gin.addr"),
			"/api",
			scene.BuildApps[sgin.GinApplication](builders),
			sgin.WithRecovery(),
			sgin.WithLogger(nil),
			sgin.WithCors(),
		),
	)
	if err := engine.Run(); err != nil {
		registry.Logger.Errorf("engine error: %s", err)
	}
}
