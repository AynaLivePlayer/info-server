package repository

import (
	"github.com/AynaLivePlayer/miaosic"
	neateasePvdr "github.com/AynaLivePlayer/miaosic/providers/netease"
	"infoserver/lyric"
)

type neteaseSearchRepo struct {
	provider miaosic.MediaProvider
}

func (n *neteaseSearchRepo) RepoImplName() string {
	return "lyric.repository.LyricSearchRepository.netease"
}

func (n *neteaseSearchRepo) Status() error {
	return nil
}

func (n *neteaseSearchRepo) Search(title string, artist string) (result lyric.Song, err error) {
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
	result.Title = song.Title
	result.Artist = song.Artist
	result.Lyrics = make([]lyric.Lyric, 0)
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

func NewNeteaseSearchRepo() lyric.LyricSearchRepository {
	return &neteaseSearchRepo{
		provider: neateasePvdr.NewNetease(),
	}
}
