package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"infoserver/bilisrv"
	"net/http"
)

type biliLiveDanmuInfoServiceImpl struct {
	repo bilisrv.BiliCredentialRepo
}

func (srv *biliLiveDanmuInfoServiceImpl) SrvImplName() string {
	return "bilisrv.BiliLiveDanmuInfoService"
}

var _ bilisrv.BiliLiveDanmuInfoService = (*biliLiveDanmuInfoServiceImpl)(nil)

func NewBiliLiveDanmuInfoService(repo bilisrv.BiliCredentialRepo) bilisrv.BiliLiveDanmuInfoService {
	return &biliLiveDanmuInfoServiceImpl{
		repo: repo,
	}
}

func (srv *biliLiveDanmuInfoServiceImpl) generateAuthBody(uid int, roomId int, token string) string {
	data := map[string]interface{}{
		"uid":    uid,
		"roomid": roomId,
		// protover  = 3 unknown encryption
		"protover": 2,
		"platform": "web",
		"type":     2,
		"key":      token,
	}
	val, _ := json.Marshal(data)
	return string(val)
}

func (srv *biliLiveDanmuInfoServiceImpl) GetDanmuInfo(roomId int) (bilisrv.BiliLiveDanmuInfo, error) {
	acc, ok := srv.repo.GetRandom()
	if !ok {
		return bilisrv.BiliLiveDanmuInfo{}, bilisrv.ErrFailToGetDanmuInfo.WithDetailStr("no credential available")
	}

	resp, err := resty.New().R().
		SetCookies([]*http.Cookie{
			{
				Name:  "SESSDATA",
				Value: acc.SessionData,
			},
			{
				Name:  "bili_jct",
				Value: acc.BilibiliJCT,
			},
		}).
		SetQueryParam("id", cast.ToString(roomId)).
		SetQueryParam("type", "0").
		Get("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo")
	if err != nil {
		return bilisrv.BiliLiveDanmuInfo{}, bilisrv.ErrFailToGetDanmuInfo.WithDetail(err)
	}
	result := gjson.Parse(string(resp.Body()))
	token := result.Get("data.token").String()
	if result.Get("code").Int() != 0 || token == "" {
		return bilisrv.BiliLiveDanmuInfo{}, bilisrv.ErrFailToGetDanmuInfo.WithDetailStr(result.Get("message").String())
	}
	var info bilisrv.BiliLiveDanmuInfo
	if err := json.Unmarshal([]byte(result.Get("data.host_list").String()), &info.HostList); err != nil {
		return bilisrv.BiliLiveDanmuInfo{}, bilisrv.ErrFailToGetDanmuInfo.WithDetail(err)
	}
	info.AuthBody = srv.generateAuthBody(acc.UID, roomId, token)
	info.WssLink = make([]string, len(info.HostList))
	for i, host := range info.HostList {
		info.WssLink[i] = fmt.Sprintf("wss://%s:%d/sub", host.Host, host.WssPort)
	}
	info.UID = acc.UID
	info.Token = token
	return info, nil
}
