package updater

import "scene-service/meowsic"

func NewNetease() meowsic.ISourceStatusUpdater {
	return &universalStatusUpdater{
		provider:      "netease",
		hasVipApi:     false,
		searchKeyword: "染 reol",
		infoMeta:      "33516503",
		infoTitle:     "染",
		fileMeta:      "33516503",
		vipFileMeta:   "33516503",
		lyricMeta:     "33516503",
		playlistMeta: []string{
			"https://music.163.com/#/playlist?id=2015397679",
		},
	}
}
