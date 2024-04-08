package blivedm

import (
	openblive "github.com/aynakeya/open-bilibili-live"
	"github.com/rhine-tech/scene"
)

type OpenBLiveApiService interface {
	scene.Service
	AppStart(code string, appId int64) (*openblive.AppStartResult, *openblive.PublicError)
	AppEnd(appId int64, gameId string) *openblive.PublicError
	HearBeat(gameId string) *openblive.PublicError
}
