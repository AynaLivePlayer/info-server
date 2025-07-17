package version

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/model"
)

var Lens scene.ModuleName = "version"

type VersionRepository interface {
	scene.Named
	GetLatest() (VersionInfo, error)
	ListVersions(offset, limit int64) (model.PaginationResult[VersionInfo], error)
	GetVersion(version Version) (VersionInfo, error)
	UpsertVersion(version VersionInfo) error
}

type VersionService interface {
	scene.Service
	CheckUpdate(clientVersion Version) (VersionInfo, bool)
	GetLatest() (VersionInfo, error)
	ListVersions(offset, limit int64) (model.PaginationResult[VersionInfo], error)
	GetVersion(version Version) (VersionInfo, error)
	UpsertVersion(version VersionInfo) error
}
