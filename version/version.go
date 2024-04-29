package version

import "github.com/rhine-tech/scene/model"

type VersionRepository interface {
	ListVersions(offset, limit int64) (model.PaginationResult[VersionInfo], error)
	GetVersion(version Version) (VersionInfo, error)
	UpsertVersion(version VersionInfo) error
}
