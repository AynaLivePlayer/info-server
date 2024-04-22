package service

import (
	"fmt"
	webApi "github.com/AynaLivePlayer/blivedm-go/api"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/infrastructure/asynctask"
	"github.com/rhine-tech/scene/infrastructure/logger"
	"infoserver/blivedm"
)

type webDanmuSingleCredImpl struct {
	logger   logger.ILogger           `aperture:""`
	td       asynctask.TaskDispatcher `aperture:""`
	biliJCT  string
	sessData string
	uid      int
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
	return w.uid, result, err
}
