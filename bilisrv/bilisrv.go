package bilisrv

import (
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

type BiliLiveDanmuInfoService interface {
	scene.Service
	GetDanmuInfo(roomId int) (BiliLiveDanmuInfo, error)
}

type BiliCredential struct {
	UID           int    `json:"uid" json:"uid"`
	AccountNumber string `json:"account_number" bson:"account_number"`
	Password      string `json:"password" bson:"password"`
	SessionData   string `json:"session_data" bson:"session_data"`
	BilibiliJCT   string `json:"bilibili_jct" bson:"bilibili_jct"`
}

type BiliCredentialRepo interface {
	scene.Repository
	Size() int
	Get(uid int) (BiliCredential, bool)
	GetRandom() (BiliCredential, bool)
	Upsert(credential BiliCredential) error
	Delete(uid int) error
	ListCredentials(offset int64, limit int64) (result []BiliCredential, total int)
}
