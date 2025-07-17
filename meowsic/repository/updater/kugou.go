package updater

import (
	"scene-service/meowsic"
)

func NewKugou() meowsic.ISourceStatusUpdater {
	return &universalStatusUpdater{
		provider:      "kugou",
		hasVipApi:     true,
		searchKeyword: "心似烟火",
		infoMeta:      "24aae0ef48311770043044ab2376a8db",
		infoTitle:     "苏星婕 - 心似烟火",
		fileMeta:      "b9a6c3eee00a7df6ff389ad383be5cb1",
		vipFileMeta:   "24aae0ef48311770043044ab2376a8db",
		lyricMeta:     "24aae0ef48311770043044ab2376a8db",
		playlistMeta: []string{
			"gcid_3zfcfgjcz31z06d",
		},
	}
}

func NewKugouInstrumental() meowsic.ISourceStatusUpdater {
	return &universalStatusUpdater{
		provider:      "kugou-instr",
		hasVipApi:     true,
		searchKeyword: "心似烟火",
		infoMeta:      "24aae0ef48311770043044ab2376a8db",
		infoTitle:     "苏星婕 - 心似烟火",
		fileMeta:      "b9a6c3eee00a7df6ff389ad383be5cb1",
		vipFileMeta:   "24aae0ef48311770043044ab2376a8db",
		lyricMeta:     "24aae0ef48311770043044ab2376a8db",
		playlistMeta: []string{
			"gcid_3zfcfgjcz31z06d",
		},
	}
}
