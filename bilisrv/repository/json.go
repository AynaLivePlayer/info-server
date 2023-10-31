package repository

import (
	"github.com/rhine-tech/scene/drivers/repos"
	"github.com/rhine-tech/scene/lens/infrastructure/logger"
	"github.com/rhine-tech/scene/model"
	"github.com/spf13/cast"
	"infoserver/bilisrv"
	"math/rand"
)

type jsonRepo struct {
	internal *repos.CommonJsonRepo[map[string]*bilisrv.BiliCredential]
	logger   logger.ILogger `aperture:""`
}

func (j *jsonRepo) Dispose() error {
	return j.internal.Dispose()
}

func (j *jsonRepo) Setup() error {
	j.logger = j.logger.WithPrefix(j.RepoImplName())
	err := j.internal.Setup()
	if err != nil {
		j.logger.Errorf("setup json repo failed: %s", err)
		return bilisrv.ErrDBConnectionError.WithDetailStr("filename: '" + j.internal.Cfg.Path + "'")
	} else {
		j.logger.Infof("success load %d credentials from json file", j.Size())
	}
	return err

}

func (j *jsonRepo) RepoImplName() string {
	return "bilisrv.repository.json"
}

func (j *jsonRepo) Status() error {
	return nil
}

func NewBiliCredentialJsonRepo(cfg model.FileConfig) bilisrv.BiliCredentialRepo {
	return &jsonRepo{
		internal: repos.NewJsonRepository[map[string]*bilisrv.BiliCredential](cfg),
	}
}

func (j *jsonRepo) Size() int {
	return len(j.internal.Data)
}

func (j *jsonRepo) Get(uid int) (bilisrv.BiliCredential, bool) {
	c, ok := j.internal.Data[cast.ToString(uid)]
	if ok {
		return *c, true
	}
	return bilisrv.BiliCredential{}, false
}

func (j *jsonRepo) GetRandom() (bilisrv.BiliCredential, bool) {
	keys := make([]string, 0, len(j.internal.Data))
	for k := range j.internal.Data {
		keys = append(keys, k)
	}

	if len(keys) == 0 {
		return bilisrv.BiliCredential{}, false
	}

	randomKey := keys[rand.Intn(len(keys))]

	return *j.internal.Data[randomKey], true
}

func (j *jsonRepo) Upsert(credential bilisrv.BiliCredential) error {
	j.internal.Data[cast.ToString(credential.UID)] = &credential
	return nil
}

func (j *jsonRepo) Delete(uid int) error {
	uidStr := cast.ToString(uid)
	if _, ok := j.internal.Data[cast.ToString(uid)]; !ok {
		return bilisrv.ErrNoSuchCredential.WithDetailStr(uidStr)
	}
	delete(j.internal.Data, cast.ToString(uid))
	return nil
}

func (j *jsonRepo) ListCredentials(offset int64, limit int64) (result []bilisrv.BiliCredential, total int) {
	total = len(j.internal.Data)
	if offset < 0 || limit <= 0 || offset >= int64(total) {
		return make([]bilisrv.BiliCredential, 0), total
	}

	// Calculate start and end indexes for the loop
	start := offset
	end := offset + limit
	if end > int64(total) {
		end = int64(total)
	}
	idx := 0
	for _, val := range j.internal.Data {
		if int64(idx) >= start && int64(idx) < end {
			result = append(result, *val)
		}
	}

	return result, len(result)
}
