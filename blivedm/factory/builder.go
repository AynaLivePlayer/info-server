package factory

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/registry"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/blivedm"
	"infoserver/blivedm/delivery"
	"infoserver/blivedm/repository/connlog"
	"infoserver/blivedm/service"
)

type App struct {
	scene.ModuleFactory
	OpenBiliLiveAccessKey    string
	OpenBiliLiveAccessSecret string
	BilibiliJCT              string
	BilibiliSessData         string
	BilibiliBuvid3           string
}

func (b App) Default() App {
	return App{
		OpenBiliLiveAccessKey:    registry.Config.GetString("blivedm.openblive.access_key"),
		OpenBiliLiveAccessSecret: registry.Config.GetString("blivedm.openblive.access_secret"),
		BilibiliJCT:              registry.Config.GetString("blivedm.bilibili.jct"),
		BilibiliSessData:         registry.Config.GetString("blivedm.bilibili.sessdata"),
		BilibiliBuvid3:           registry.Config.GetString("blivedm.bilibili.buvid3"),
	}
}

func (b App) Init() scene.LensInit {
	return func() {
		registry.Register(connlog.GormRepo(nil))
	}
}

func (b App) Apps() []any {
	return []any{
		func() sgin.GinApplication {
			return delivery.GinApp(
				registry.Load[blivedm.WebDanmuService](
					service.NewWebDanmuServiceSingleCredential(b.BilibiliJCT, b.BilibiliSessData, b.BilibiliBuvid3)),
				registry.Load[blivedm.OpenBLiveApiService](
					service.NewOpenBLiveApiService(b.OpenBiliLiveAccessKey, b.OpenBiliLiveAccessSecret)),
				registry.Load[blivedm.ConnectionLogService](service.ConnectionLogService(nil)),
			)
		},
	}
}
