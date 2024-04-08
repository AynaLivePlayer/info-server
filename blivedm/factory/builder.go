package factory

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/registry"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/blivedm"
	"infoserver/blivedm/delivery"
	"infoserver/blivedm/service"
)

type App struct {
	scene.ModuleFactory
	OpenBiliLiveAccessKey    string
	OpenBiliLiveAccessSecret string
	BilibiliJCT              string
	BilibiliSessData         string
}

func (b App) Default() App {
	return App{
		OpenBiliLiveAccessKey:    registry.Config.GetString("blivedm.openblive.access_key"),
		OpenBiliLiveAccessSecret: registry.Config.GetString("blivedm.openblive.access_secret"),
		BilibiliJCT:              registry.Config.GetString("blivedm.bilibili.jct"),
		BilibiliSessData:         registry.Config.GetString("blivedm.bilibili.sessdata"),
	}
}

func (b App) Init() scene.LensInit {
	return func() {

	}
}

func (b App) Apps() []any {
	return []any{

		func() sgin.GinApplication {
			openblive := registry.Load[blivedm.OpenBLiveApiService](
				service.NewOpenBLiveApiService(b.OpenBiliLiveAccessKey, b.OpenBiliLiveAccessSecret))
			dmSrv := registry.Load[blivedm.WebDanmuService](
				service.NewWebDanmuServiceSingleCredential(b.BilibiliJCT, b.BilibiliSessData))
			return registry.Load(delivery.NewGinApp(dmSrv, openblive))
		},
	}
}
