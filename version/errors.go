package version

import "github.com/rhine-tech/scene/errcode"

var _eg = errcode.NewErrorGroup(13, "version")

var (
	ErrFailedToGetLatest    = _eg.CreateError(1, "failed to get latest version")
	ErrFailedToListVersions = _eg.CreateError(2, "failed to list versions")
	ErrFailedToGetVersion   = _eg.CreateError(3, "failed to get version")
	ErrFailToUpdateVersion  = _eg.CreateError(4, "failed to update version")
)
