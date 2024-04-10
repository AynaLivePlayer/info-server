package main

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/engines"
	asynctask "github.com/rhine-tech/scene/lens/infrastructure/asynctask/factory"
	config "github.com/rhine-tech/scene/lens/infrastructure/config/factory"
	logger "github.com/rhine-tech/scene/lens/infrastructure/logger/factory"
	"github.com/rhine-tech/scene/registry"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	blivedm "infoserver/blivedm/factory"
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
		asynctask.Ants{},
		blivedm.App{}.Default(),
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
