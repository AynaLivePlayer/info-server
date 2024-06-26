package streamerstat

import "github.com/rhine-tech/scene/errcode"

var _eg = errcode.NewErrorGroup(20, Lens.String())

var (
	ErrFailToUpdateStatus    = _eg.CreateError(1, "fail to update status")
	ErrStatusNotFound        = _eg.CreateError(2, "status not found")
	ErrGetStatusUnknownError = _eg.CreateError(3, "get status error")
)
