package lyric

import "github.com/rhine-tech/scene/errcode"

var _eg = errcode.NewErrorGroup(12, "lyric")

var (
	ErrCantFindLyric = _eg.CreateError(1, "can't find lyric")
)
