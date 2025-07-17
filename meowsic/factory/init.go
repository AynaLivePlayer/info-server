package factory

import (
	_ "github.com/AynaLivePlayer/miaosic"
	_ "github.com/AynaLivePlayer/miaosic/providers/bilivideo"
	"github.com/AynaLivePlayer/miaosic/providers/kugou"
	_ "github.com/AynaLivePlayer/miaosic/providers/kugou"
	_ "github.com/AynaLivePlayer/miaosic/providers/kuwo"
	_ "github.com/AynaLivePlayer/miaosic/providers/local"
	_ "github.com/AynaLivePlayer/miaosic/providers/netease"
	_ "github.com/AynaLivePlayer/miaosic/providers/qq"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/registry"
	"scene-service/meowsic"
	"scene-service/meowsic/delivery"
	"scene-service/meowsic/repository/sessionstore"
	"scene-service/meowsic/repository/statusstore"
	"scene-service/meowsic/repository/updater"
	"scene-service/meowsic/service"
	"scene-service/onebot"
)

func init() {
	kugou.UseInstrumental()
}

type App struct {
	scene.ModuleFactory
}

func (b App) Init() scene.LensInit {
	return func() {
		_ = registry.Register(sessionstore.NewGorm())
		_ = registry.Register(statusstore.NewGorm())
		_ = registry.Register[meowsic.IMeowsicService](
			service.NewMeowsicService(
				registry.Load(updater.NewKugou()),
				registry.Load(updater.NewKugouInstrumental()),
				registry.Load(updater.NewNetease()),
				registry.Load(updater.NewBiliVideo()),
				registry.Load(updater.NewKuwo()),
				registry.Load(updater.NewQQ())),
		)
	}
}

func (b App) Apps() []any {
	return []any{
		delivery.GinApp,
		func() onebot.ChatBotPlugin {
			return registry.Load(delivery.NewChatBotPlugin())
		},
	}
}
