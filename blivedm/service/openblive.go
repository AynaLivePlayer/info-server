package service

import (
	"github.com/aynakeya/open-bilibili-live"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/lens/infrastructure/logger"
	"infoserver/blivedm"
)

type openBLiveApiServiceImpl struct {
	logger       logger.ILogger `aperture:""`
	apiClient    *openblive.ApiClient
	accessKey    string
	accessSecret string
}

func (o *openBLiveApiServiceImpl) Setup() error {
	o.logger = o.logger.WithPrefix(o.SrvImplName().Identifier())
	if len(o.accessKey) < 8 || len(o.accessKey) < 8 {
		o.logger.Errorf("invalid openblive access key or access secret")
		return blivedm.ErrInvalidOpenBLiveApiService
	}
	o.logger.InfoW("initialize with", "openblive access key", o.accessKey[:4], "openblive access secret", o.accessSecret[:4])
	return nil
}

func NewOpenBLiveApiService(accessKey, accessSecret string) blivedm.OpenBLiveApiService {
	return &openBLiveApiServiceImpl{
		apiClient:    openblive.NewApiClient(accessKey, accessSecret),
		accessKey:    accessKey,
		accessSecret: accessSecret,
	}
}

func (o *openBLiveApiServiceImpl) SrvImplName() scene.ImplName {
	return scene.NewSrvImplNameNoVer(blivedm.ModuleName, "OpenBLiveApiService")
}

func (o *openBLiveApiServiceImpl) AppStart(code string, appId int64) (*openblive.AppStartResult, *openblive.PublicError) {
	return o.apiClient.AppStart(code, appId)
}

func (o *openBLiveApiServiceImpl) AppEnd(appId int64, gameId string) *openblive.PublicError {
	return o.apiClient.AppEnd(appId, gameId)
}

func (o *openBLiveApiServiceImpl) HearBeat(gameId string) *openblive.PublicError {
	return o.apiClient.HearBeat(gameId)
}
