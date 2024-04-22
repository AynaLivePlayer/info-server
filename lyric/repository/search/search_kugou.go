package search

//import (
//	"encoding/base64"
//	"github.com/aynakeya/deepcolor"
//	"github.com/aynakeya/deepcolor/dphttp"
//	"github.com/tidwall/gjson"
//	"infoserver/lyric"
//)
//
//type kugouLrcCandidate struct {
//	ID        string `json:"id"`
//	AccessKey string `json:"access_key"`
//	Title     string `json:"song"`
//	Artist    string `json:"singer"`
//	Language  string `json:"language"`
//}
//
//type kugouSearchRepo struct {
//	searchApi dphttp.ApiResultFunc[string, []kugouLrcCandidate]
//	lyricApi  dphttp.ApiResultFunc[kugouLrcCandidate, string]
//}
//
//func NewKugouLyricSearch() lyric.LyricProvider {
//	kg := &kugouSearchRepo{}
//	kg.searchApi = deepcolor.CreateApiResultFunc(
//		deepcolor.NewGetRequestFuncWithSingleQuery(
//			"http://lyrics.kugou.com/search?ver=1&man=yes&client=pc&hash=",
//			"keyword",
//			nil),
//		deepcolor.ParserGJson,
//		func(result *gjson.Result, container *[]kugouLrcCandidate) error {
//			for _, v := range result.Get("candidates").Array() {
//				c := kugouLrcCandidate{
//					ID:        v.Get("id").String(),
//					AccessKey: v.Get("accesskey").String(),
//					Title:     v.Get("song").String(),
//					Artist:    v.Get("singer").String(),
//					Language:  v.Get("language").String(),
//				}
//				if c.Language == "" {
//					c.Language = "default"
//				}
//				*container = append(*container, c)
//			}
//			return nil
//		})
//	kg.lyricApi = deepcolor.CreateApiResultFunc(
//		func(params kugouLrcCandidate) (*dphttp.Request, error) {
//			return deepcolor.NewGetRequestWithQuery(
//				"http://lyrics.kugou.com/download?ver=1&client=pc&fmt=lrc&charset=utf8",
//				map[string]any{
//					"accesskey": params.AccessKey,
//					"id":        params.ID,
//				}, nil)
//		},
//		deepcolor.ParserGJson,
//		func(result *gjson.Result, container *string) error {
//			data, err := base64.StdEncoding.DecodeString(result.Get("content").String())
//			if err != nil {
//				return err
//			}
//			*container = string(data)
//			return nil
//		})
//	return kg
//}
//
//func (k *kugouSearchRepo) RepoImplName() string {
//	return "lyric.repository.LyricProvider.kugou"
//}
//
//func (k *kugouSearchRepo) Status() error {
//	return nil
//}
//
//func (k *kugouSearchRepo) Search(title string, artist string) (result lyric.Song, err error) {
//	songs, err := k.searchApi(title + " " + artist)
//	if err != nil {
//		return result, err
//	}
//	if len(songs) == 0 {
//		return result, lyric.ErrCantFindLyric
//	}
//	var song kugouLrcCandidate
//	for _, v := range songs {
//		if v.Title == title {
//			song = v
//			break
//		}
//	}
//	if song.Title == "" {
//		return result, lyric.ErrCantFindLyric
//	}
//	lyrics, err := k.lyricApi(song)
//	if err != nil {
//		return result, err
//	}
//	result.Title = song.Title
//	result.Artist = song.Artist
//	result.Lyrics = make([]lyric.Lyric, 0)
//	result.Lyrics = append(result.Lyrics, lyric.Lyric{
//		Lang:  song.Language,
//		Lyric: lyrics,
//	})
//	if len(result.Lyrics) == 0 {
//		return result, lyric.ErrCantFindLyric
//	}
//	return result, nil
//}
