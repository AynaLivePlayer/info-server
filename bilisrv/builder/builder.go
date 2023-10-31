package builder

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/lens/infrastructure/logger"
	"github.com/rhine-tech/scene/model"
	"github.com/rhine-tech/scene/registry"
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/bilisrv"

	"infoserver/bilisrv/delivery"
	"infoserver/bilisrv/repository"
	"infoserver/bilisrv/service"
)

func CreateApp(
	logger logger.ILogger,
	cfg model.DatabaseConfig) sgin.GinApplication {
	repo := registry.Register(repository.NewBiliCredentialMongoRepo(cfg))
	srv2 := registry.Register(service.NewBiliLiveDanmuInfoService(repo))
	return delivery.NewGinApp(logger, srv2)
}

func InitApp() sgin.GinApplication {
	return CreateApp(
		registry.AcquireSingleton(logger.ILogger(nil)),
		*registry.AcquireSingleton(&model.DatabaseConfig{}))
}

type JsonDB struct {
	scene.Builder
}

func (b JsonDB) Init() scene.LensInit {
	return func() {
		repo := registry.Register(
			repository.NewBiliCredentialJsonRepo(model.NewFileConfig(registry.Config.GetString("bilisrv.jsondb"))))
		registry.Register(service.NewBiliLiveDanmuInfoService(repo))
	}
}

func (b JsonDB) Apps() []any {
	return []any{
		func() sgin.GinApplication {
			return delivery.NewGinApp(
				registry.Use(logger.ILogger(nil)),
				registry.Use(bilisrv.BiliLiveDanmuInfoService(nil)))
		},
	}
}

type Builder struct {
	scene.Builder
}

func (b Builder) Apps() []any {
	return []any{
		InitApp,
	}
}
