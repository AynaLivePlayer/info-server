package blivedm

import (
	webApi "github.com/AynaLivePlayer/blivedm-go/api"
	"github.com/rhine-tech/scene"
)

type BiliLiveDanmuInfo struct {
	HostList []struct {
		Host    string `json:"host"`
		Port    int    `json:"port"`
		WsPort  int    `json:"ws_port"`
		WssPort int    `json:"wss_port"`
	} `json:"host_list"`
	Token    string   `json:"token"`
	UID      int      `json:"uid"`
	WssLink  []string `json:"wss_link"`
	AuthBody string   `json:"auth_body"`
}

type WebDanmuService interface {
	scene.Service
	GetDanmuInfo(roomID int) (int, *webApi.DanmuInfo, error)
	// GetDanmuInfoCompatible don't remember what fuck is this for,
	// but keep it as is
	GetDanmuInfoCompatible(roomID int) (BiliLiveDanmuInfo, error)
}
