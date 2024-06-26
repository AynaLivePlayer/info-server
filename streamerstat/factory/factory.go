package factory

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/registry"
	"infoserver/streamerstat/repository/storage"
	"infoserver/streamerstat/repository/updaters"
	"infoserver/streamerstat/service"
)

type App struct {
	scene.ModuleFactory
}

func (b App) Init() scene.LensInit {
	return func() {
		repo := registry.Load(storage.GormRepository(nil))
		registry.Register(service.StreamerStatsService(
			repo, new(updaters.Bilibili),
		))
	}
}
