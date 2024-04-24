package service

import (
	"github.com/sahilm/fuzzy"
	"infoserver/lyric"
	"sort"
	"strings"
)

type rankScore struct {
	media *lyric.Song
	score int
}

func rankMedia(keyword string, songs []lyric.Song) []lyric.Song {
	patterns := strings.Split(keyword, " ")
	data := make([]*rankScore, 0)

	for i, _ := range songs {
		data = append(data, &rankScore{
			media: &songs[i],
			score: 0,
		})
	}

	for _, pattern := range patterns {
		pattern = strings.ToLower(pattern)
		dataStr := make([]string, 0)
		for _, d := range data {
			dataStr = append(dataStr, strings.ToLower(d.media.Title))
		}
		for _, match := range fuzzy.Find(pattern, dataStr) {
			data[match.Index].score += match.Score
		}

		dataStr = make([]string, 0)
		for _, d := range data {
			dataStr = append(dataStr, strings.ToLower(strings.Join(d.media.Artist, ",")))
		}
		for _, match := range fuzzy.Find(pattern, dataStr) {
			data[match.Index].score += match.Score
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].score > data[j].score
	})

	result := make([]lyric.Song, 0)
	for _, d := range data {
		if d.score > 0 {
			result = append(result, *d.media)
		}
	}
	return result
}
