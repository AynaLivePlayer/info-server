package blivedm

import (
	webApi "github.com/AynaLivePlayer/blivedm-go/api"
	"github.com/rhine-tech/scene"
)

type WebDanmuService interface {
	scene.Service
	GetDanmuInfo(roomID int) (int, *webApi.DanmuInfo, error)
}
