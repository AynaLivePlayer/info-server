package storage

import (
	"fmt"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/composition/database"
	"gorm.io/gorm"
	"infoserver/lyric"
)

type mysqlStorageImpl struct {
	gorm database.Gorm `aperture:""`
}

func MysqlImpl(gorm database.Gorm) lyric.LyricStorageRepository {
	return &mysqlStorageImpl{gorm: gorm}
}

func (m *mysqlStorageImpl) Setup() error {
	err := m.gorm.RegisterModel(&tableArtist{}, &tableLyric{}, &tableSong{})
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlStorageImpl) RepoImplName() scene.ImplName {
	return scene.NewModuleImplName("lyric", "LyricStorageRepository", "mysql")
}

func (m *mysqlStorageImpl) GetLyric(title string, artist string) (result []lyric.Song, err error) {
	var songs []tableSong
	err = m.gorm.DB().
		Preload("Artists").
		Preload("Lyrics").
		Distinct("lyric_songs.song_id").
		Select("lyric_songs.song_id, lyric_songs.title").
		Joins("JOIN lyric_song_artists ON lyric_song_artists.table_song_song_id = lyric_songs.song_id").
		Joins("JOIN lyric_artists ON lyric_artists.artist_id = lyric_song_artists.table_artist_artist_id").
		Joins("JOIN lyric_lyrics ON lyric_songs.song_id = lyric_lyrics.song_id").
		Where("(LOWER(lyric_songs.title) = LOWER(?) OR '' = ?) AND ('' = ? OR LOWER(lyric_artists.name) = LOWER(?))", title, title, artist, artist).
		Find(&songs).Error

	if err != nil {
		return nil, err
	}

	result = make([]lyric.Song, len(songs))
	for i, song := range songs {
		fmt.Println(song.SongID, song.Title, song.Artists)
		result[i] = lyric.Song{
			Title: song.Title,
		}
		for _, artist := range song.Artists {
			result[i].Artist = append(result[i].Artist, artist.Name)
		}
		for _, lrc := range song.Lyrics {
			result[i].Lyrics = append(result[i].Lyrics, lyric.Lyric{Lang: lrc.Lang, Lyric: lrc.LyricText})
		}
	}

	return result, nil
}

func (m *mysqlStorageImpl) Add(songs []lyric.Song) (err error) {
	// Transactional insertion to handle multiple songs and related data
	err = m.gorm.DB().Transaction(func(tx *gorm.DB) error {
		for _, song := range songs {
			tSong := tableSong{
				Title:   song.Title,
				Artists: make([]tableArtist, 0),
				Lyrics:  make([]tableLyric, 0),
			}
			for _, artist := range song.Artist {
				tSong.Artists = append(tSong.Artists, tableArtist{Name: artist})
			}
			for _, lrc := range song.Lyrics {
				tSong.Lyrics = append(tSong.Lyrics, tableLyric{Lang: lrc.Lang, LyricText: lrc.Lyric})
			}
			if err := tx.Create(&tSong).Error; err != nil {
				return err // Rollback the transaction on error
			}
		}
		return nil
	})
	return err
}
