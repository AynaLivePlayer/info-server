package service

import (
	"github.com/stretchr/testify/require"
	"scene-service/bilisrv"
	"testing"
)

type dummyBiliCredRepo struct{}

func (d dummyBiliCredRepo) RepoImplName() string {
	//TODO implement me
	panic("implement me")
}

func (d dummyBiliCredRepo) Status() error {
	//TODO implement me
	panic("implement me")
}

func (d dummyBiliCredRepo) Size() int {
	//TODO implement me
	panic("implement me")
}

func (d dummyBiliCredRepo) Upsert(credential bilisrv.BiliCredential) error {
	//TODO implement me
	panic("implement me")
}

func (d dummyBiliCredRepo) Delete(uid int) error {
	//TODO implement me
	panic("implement me")
}

func (d dummyBiliCredRepo) ListCredentials(offset int64, limit int64) (result []bilisrv.BiliCredential, total int) {
	//TODO implement me
	panic("implement me")
}

func (d dummyBiliCredRepo) Get(uid int) (bilisrv.BiliCredential, bool) {
	return bilisrv.BiliCredential{
		UID:         549644123,
		SessionData: "",
		BilibiliJCT: "",
	}, true
}

func (d dummyBiliCredRepo) GetRandom() (bilisrv.BiliCredential, bool) {
	return bilisrv.BiliCredential{
		UID:         549644123,
		SessionData: "",
		BilibiliJCT: "",
	}, true
}

func TestBiliLiveDanmuInfoServiceImpl_GetDanmuInfo(t *testing.T) {
	srv := NewBiliLiveDanmuInfoService(dummyBiliCredRepo{})
	danmuInfo, err := srv.GetDanmuInfo(7777)
	require.Nil(t, err)
	require.NotEmpty(t, danmuInfo.HostList)
}
