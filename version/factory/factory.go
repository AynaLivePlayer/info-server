package factory

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/registry"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/version/delivery"
	"infoserver/version/repository"
	"infoserver/version/service"
)

type App struct {
}

func (d App) Init() scene.LensInit {
	return func() {
	}
}

func (d App) Apps() []any {
	return []any{
		func() sgin.GinApplication {
			repo := registry.Load(repository.VersionRepository(nil))
			srv := registry.Load(service.VersionService(repo))
			return registry.Load(delivery.NewGinApp(srv))
		},
	}
}
