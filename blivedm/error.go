package blivedm

import "github.com/rhine-tech/scene/errcode"

var _eg = errcode.WithErrGroup(15)

var (
	ErrInvalidOpenBLiveApiService = _eg.CreateError(1, "fail to create open blive api service")
)
