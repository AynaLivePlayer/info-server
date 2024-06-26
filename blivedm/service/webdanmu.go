package service

import (
	"fmt"
	webApi "github.com/AynaLivePlayer/blivedm-go/api"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/infrastructure/asynctask"
	"github.com/rhine-tech/scene/infrastructure/logger"
	"infoserver/blivedm"
	"infoserver/blivedm/pkg/dmpacket"
	"infoserver/streamerstat"
	"strconv"
	"time"
)

type webDanmuSingleCredImpl struct {
	logger     logger.ILogger                     `aperture:""`
	td         asynctask.TaskDispatcher           `aperture:""`
	logRepo    blivedm.ConnectionLogRepository    `aperture:""`
	streamStat streamerstat.IStreamerStatsService `aperture:""`
	biliJCT    string
	sessData   string
	uid        int
}

func (w *webDanmuSingleCredImpl) Setup() error {
	w.logger = w.logger.WithPrefix(w.SrvImplName().Identifier())
	if len(w.biliJCT) < 8 || len(w.sessData) < 8 {
		w.logger.Errorf("invalid biliJCT or sessData")
		return blivedm.ErrInvalidWebDanmuService
	}
	w.td.Run(func() error {
		w.updateUid()
		return nil
	})
	return nil
}

func NewWebDanmuServiceSingleCredential(biliJCT, sessData string) blivedm.WebDanmuService {
	return &webDanmuSingleCredImpl{
		biliJCT:  biliJCT,
		sessData: sessData,
		uid:      -1,
	}
}

func (w *webDanmuSingleCredImpl) SrvImplName() scene.ImplName {
	return blivedm.Lens.ImplName("WebDanmuService", "SingleCredential")
}

func (w *webDanmuSingleCredImpl) cookie() string {
	return fmt.Sprintf("bili_jct=%s; SESSDATA=%s", w.biliJCT, w.sessData)

}

func (w *webDanmuSingleCredImpl) updateUid() {
	w.logger.InfoW("try to get uid", "biliJCT", w.biliJCT[:4], "sessData", w.sessData[:4])
	uid, err := webApi.GetUid(w.cookie())
	if err != nil {
		w.logger.ErrorW("failed to get uid", "err", err)
		w.uid = 0
		return
	}
	w.uid = uid
	w.logger.InfoW("uid updated", "uid", w.uid)
}

func (w *webDanmuSingleCredImpl) GetDanmuInfo(roomID int) (int, *webApi.DanmuInfo, error) {
	// uid not initialized
	if w.uid == -1 {
		w.logger.InfoW("uid not initialized, try to initialize uid")
		w.updateUid()
	}
	result, err := webApi.GetDanmuInfo(roomID, w.cookie())
	tn := time.Now().Unix()
	w.td.Run(func() error {
		_ = w.logRepo.AddEntry(roomID, "webdanmu", tn)
		_, _ = w.streamStat.UpdateStatus("bilibili", strconv.Itoa(roomID))
		return nil
	})
	return w.uid, result, err
}

func (w *webDanmuSingleCredImpl) GetDanmuInfoCompatible(roomID int) (info blivedm.BiliLiveDanmuInfo, err error) {
	uid, result, err := w.GetDanmuInfo(roomID)
	if err != nil {
		return
	}
	info.AuthBody = dmpacket.GenerateAuthBody(uid, roomID, result.Data.Token)
	for _, host := range result.Data.HostList {
		info.HostList = append(info.HostList, struct {
			Host    string `json:"host"`
			Port    int    `json:"port"`
			WsPort  int    `json:"ws_port"`
			WssPort int    `json:"wss_port"`
		}{
			Host:    host.Host,
			Port:    host.Port,
			WssPort: host.WssPort,
			WsPort:  host.WsPort,
		})
		info.WssLink = append(info.WssLink, fmt.Sprintf("wss://%s:%d/sub", host.Host, host.WssPort))
	}
	info.UID = w.uid
	info.Token = result.Data.Token
	return info, nil
}
