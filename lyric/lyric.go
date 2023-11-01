package lyric

import (
	"fmt"
	"github.com/rhine-tech/scene"
)

type Song struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Lyrics []Lyric `json:"lyrics"`
}

func (s *Song) String() string {
	return fmt.Sprintf("<lyric.Song %s %s>", s.Title, s.Artist)
}

type Lyric struct {
	Lang  string `json:"lang"`
	Lyric string `json:"lyric"`
}

type LyricStorageRepository interface {
	scene.Repository
	Search(title string, artist string) (result []Song, err error)
	Add(song Song) (err error)
}

type LyricSearchRepository interface {
	scene.Repository
	Search(title string, artist string) (result Song, err error)
}

type LyricService interface {
	scene.Service
	Search(title string, artist string) (result []Song, err error)
}
