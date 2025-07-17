package search

import (
	"github.com/AynaLivePlayer/miaosic"
	"github.com/AynaLivePlayer/miaosic/providers/netease"
	"github.com/rhine-tech/scene"
	"infoserver/lyric"
	"strings"
)

type neteaseSearchRepo struct {
	provider miaosic.MediaProvider
}

func NewNeteaseProvider() lyric.LyricProvider {
	return &neteaseSearchRepo{provider: netease.NewNetease()}
}

func (n *neteaseSearchRepo) ImplName() scene.ImplName {
	return lyric.Lens.ImplName("LyricProvider", "netease")
}

func (n *neteaseSearchRepo) GetLyric(title string, artist string) (result lyric.Song, err error) {
	songs, err := n.provider.Search(title+" "+artist, 1, 5)
	if err != nil {
		return result, err
	}
	if len(songs) == 0 {
		return result, lyric.ErrCantFindLyric
	}
	var song miaosic.MediaInfo
	for _, v := range songs {
		if v.Title == title {
			song = v
			break
		}
	}
	if song.Title == "" {
		return result, lyric.ErrCantFindLyric
	}
	lyrics, err := n.provider.GetMediaLyric(song.Meta)
	if err != nil {
		return result, err
	}
	result = lyric.Song{
		Title:  song.Title,
		Artist: strings.Split(song.Artist, ","),
		Lyrics: make([]lyric.Lyric, 0),
	}
	for _, v := range lyrics {
		result.Lyrics = append(result.Lyrics, lyric.Lyric{
			Lang:  v.Lang,
			Lyric: v.String(),
		})
	}
	if len(result.Lyrics) == 0 {
		return result, lyric.ErrCantFindLyric
	}
	return result, nil
}
