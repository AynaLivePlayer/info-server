package updater

import "scene-service/meowsic"

func NewBiliVideo() meowsic.ISourceStatusUpdater {
	return &universalStatusUpdater{
		provider:      "bilibili-video",
		hasVipApi:     false,
		searchKeyword: "家有女友op",
		infoMeta:      "BV1Mx411P75Y",
		infoTitle:     "【洛天依】世末积雨云",
		fileMeta:      "BV1Mx411P75Y",
		vipFileMeta:   "BV1Mx411P75Y",
		lyricMeta:     "BV1Mx411P75Y",
		playlistMeta: []string{
			"https://space.bilibili.com/10003632/favlist?fid=729246932&ftype=create",
		},
	}
}
