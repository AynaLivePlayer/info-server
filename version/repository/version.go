package repository

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/composition/orm"
	"github.com/rhine-tech/scene/model"
	"infoserver/version"
)

var _ version.VersionRepository = &versionImpl{}

type versionImpl struct {
	db orm.Gorm `aperture:""`
}

func VersionRepository(db orm.Gorm) version.VersionRepository {
	return &versionImpl{db: db}
}

func (v *versionImpl) Setup() error {
	err := v.db.RegisterModel(&tableVersion{})
	return err
}

func (v *versionImpl) RepoImplName() scene.ImplName {
	return version.Lens.ImplNameNoVer("VersionRepository")
}

func (v *versionImpl) GetLatest() (version.VersionInfo, error) {
	var latestVersion tableVersion
	err := v.db.DB().Order("version desc").First(&latestVersion).Error
	if err != nil {
		return version.VersionInfo{}, err
	}
	return version.VersionInfo{
		Version:     version.Version(latestVersion.Version),
		Note:        latestVersion.Note,
		ReleaseTime: latestVersion.ReleaseDate,
	}, nil
}

func (v *versionImpl) ListVersions(offset, limit int64) (result model.PaginationResult[version.VersionInfo], err error) {
	var versions []tableVersion
	err = v.db.DB().Offset(int(offset)).Limit(int(limit)).Find(&versions).Error
	if err != nil {
		return result, err
	}
	var versionInfos []version.VersionInfo
	for _, ver := range versions {
		versionInfos = append(versionInfos, version.VersionInfo{
			Version:     version.Version(ver.Version),
			Note:        ver.Note,
			ReleaseTime: ver.ReleaseDate,
		})
	}
	result.Results = versionInfos
	err = v.db.DB().Model(&tableVersion{}).Count(&result.Total).Error
	result.Count = int64(len(versionInfos))
	return result, err
}

func (v *versionImpl) GetVersion(ver version.Version) (version.VersionInfo, error) {
	var verTable tableVersion
	err := v.db.DB().Where("version = ?", uint32(ver)).First(&verTable).Error
	if err != nil {
		return version.VersionInfo{}, err
	}
	return version.VersionInfo{
		Version:     ver,
		Note:        verTable.Note,
		ReleaseTime: verTable.ReleaseDate,
	}, nil
}

func (v *versionImpl) UpsertVersion(ver version.VersionInfo) error {
	verTable := tableVersion{
		Version:     uint32(ver.Version),
		Note:        ver.Note,
		ReleaseDate: ver.ReleaseTime,
	}
	err := v.db.DB().
		Where(tableVersion{Version: uint32(ver.Version)}).
		Assign(verTable).
		FirstOrCreate(&verTable).Error
	return err
}
