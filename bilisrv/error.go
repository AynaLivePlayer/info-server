package bilisrv

import "github.com/rhine-tech/scene/errcode"

var _eg = errcode.WithErrGroup(11)

var (
	ErrDBConnectionError  = _eg.CreateError(1, "fail to load credential from db")
	ErrFailToGetDanmuInfo = _eg.CreateError(2, "fail to get danmu info")
	ErrNoSuchCredential   = _eg.CreateError(3, "no such credential")
)
