package builder

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/registry"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/lyric/delivery"
	"infoserver/lyric/repository"
	"infoserver/lyric/service"
)

type DefaultBuilder struct {
}

func (d DefaultBuilder) Init() scene.LensInit {
	return func() {
		registry.Register(service.NewLyricService(
			registry.Load(repository.NewMysqlStorageRepository()),
			registry.Load(repository.NewNeteaseSearchRepo()),
		))
	}
}

func (d DefaultBuilder) Apps() []any {
	return []any{
		func() sgin.GinApplication {
			return registry.Load(delivery.NewGinApp())
		},
	}
}
