package lyric

import (
	"fmt"
	"github.com/rhine-tech/scene"
)

const Lens scene.ModuleName = "lyric"

type Song struct {
	Title  string   `json:"title"`
	Artist []string `json:"artist"`
	Lyrics []Lyric  `json:"lyrics"`
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
	GetLyric(title string, artist string) (result []Song, err error)
	Add(song []Song) (err error)
	//Search(keyword string) (result []Song, err error)
}

type LyricProvider interface {
	scene.Repository
	GetLyric(title string, artist string) (result Song, err error)
}

type LyricService interface {
	scene.Service
	GetLyric(title string, artist string) (result []Song, err error)
	//Search(keyword string) (result []Song, err error)
}
