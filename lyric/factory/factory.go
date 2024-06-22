package factory

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/registry"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/lyric/delivery"
	"infoserver/lyric/repository/search"
	"infoserver/lyric/repository/storage"
	"infoserver/lyric/service"
)

type DefaultBuilder struct {
}

func (d DefaultBuilder) Init() scene.LensInit {
	return func() {
		registry.Register(service.NewLyricService(
			registry.Load(storage.MysqlImpl(nil)),
			registry.Load(search.NewNeteaseProvider()),
			registry.Load(search.NewKuwoProvider()),
		),
		)
	}
}

func (d DefaultBuilder) Apps() []any {
	return []any{
		func() sgin.GinApplication {
			return registry.Load(delivery.NewGinApp())
		},
	}
}
