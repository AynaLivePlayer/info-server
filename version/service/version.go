package service

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/infrastructure/logger"
	"github.com/rhine-tech/scene/model"
	"infoserver/version"
)

var _ = version.VersionService(&versionSrvImpl{})

type versionSrvImpl struct {
	repo version.VersionRepository `aperture:""`
	log  logger.ILogger            `aperture:""`
}

func VersionService(repo version.VersionRepository) version.VersionService {
	return &versionSrvImpl{repo: repo}
}

func (v *versionSrvImpl) Setup() error {
	v.log = v.log.WithPrefix(v.SrvImplName().Identifier())
	return nil
}

func (v *versionSrvImpl) SrvImplName() scene.ImplName {
	return version.Lens.ImplNameNoVer("VersionService")
}

func (v *versionSrvImpl) CheckUpdate(clientVersion version.Version) (version.VersionInfo, bool) {
	ver, err := v.repo.GetLatest()
	if err != nil {
		v.log.ErrorW("fail to get latest version when checking update", "err", err)
		return version.VersionInfo{}, false
	}
	if ver.Version > clientVersion {
		return ver, true
	}
	return ver, false
}

func (v *versionSrvImpl) GetLatest() (version.VersionInfo, error) {
	ver, err := v.repo.GetLatest()
	if err != nil {
		v.log.ErrorW("fail to get latest version", "err", err)
		return version.VersionInfo{}, version.ErrFailedToGetLatest
	}
	return ver, nil
}

func (v *versionSrvImpl) ListVersions(offset, limit int64) (model.PaginationResult[version.VersionInfo], error) {
	vers, err := v.repo.ListVersions(offset, limit)
	if err != nil {
		v.log.ErrorW("fail to list versions", "err", err)
		return model.PaginationResult[version.VersionInfo]{}, version.ErrFailedToListVersions
	}
	return vers, nil
}

func (v *versionSrvImpl) GetVersion(verNum version.Version) (version.VersionInfo, error) {
	ver, err := v.repo.GetVersion(verNum)
	if err != nil {
		v.log.ErrorW("fail to get version", "err", err)
		return version.VersionInfo{}, version.ErrFailedToGetVersion.WithDetailStr(verNum.String())
	}
	return ver, nil
}

func (v *versionSrvImpl) UpsertVersion(ver version.VersionInfo) error {
	err := v.repo.UpsertVersion(ver)
	if err != nil {
		v.log.ErrorW("fail to upsert version", "err", err)
		return version.ErrFailToUpdateVersion
	}
	return nil
}
