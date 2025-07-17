package updater

import "scene-service/meowsic"

func NewKuwo() meowsic.ISourceStatusUpdater {
	return &universalStatusUpdater{
		provider:      "kuwo",
		hasVipApi:     true,
		searchKeyword: "周杰伦",
		infoMeta:      "22804772",
		infoTitle:     "霜雪千年",
		fileMeta:      "22804772",
		vipFileMeta:   "6536164",
		lyricMeta:     "22804772",
		playlistMeta: []string{
			"",
		},
	}
}
