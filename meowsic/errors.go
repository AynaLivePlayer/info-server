package meowsic

import "github.com/rhine-tech/scene/errcode"

var _eg = errcode.NewErrorGroup(25, "meowsic")

var (
	ErrSourceNotLoginable   = _eg.CreateError(1, "target source not loginable")
	ErrSourceNotFound       = _eg.CreateError(2, "target source not found")
	ErrStatusNotFound       = _eg.CreateError(3, "target status not found")
	ErrInvalidMeta          = _eg.CreateError(4, "invalid metadata")
	ErrorRestoreLoginFailed = _eg.CreateError(5, "restore login failed")
	ErrLogoutFailed         = _eg.CreateError(6, "logout failed")
	ErrMiaosicError         = _eg.CreateError(7, "miaosic error")
)
