package updater

import "scene-service/meowsic"

func NewQQ() meowsic.ISourceStatusUpdater {
	return &universalStatusUpdater{
		provider:      "qq",
		hasVipApi:     false,
		searchKeyword: "长夜雨",
		infoMeta:      "0001bGCK1lQoYf",
		infoTitle:     "长夜雨",
		fileMeta:      "0001bGCK1lQoYf",
		vipFileMeta:   "001gP4t40Q0EO0",
		lyricMeta:     "001gP4t40Q0EO0",
		playlistMeta: []string{
			"https://y.qq.com/n/ryqq/playlist/7426999757",
		},
	}
}
