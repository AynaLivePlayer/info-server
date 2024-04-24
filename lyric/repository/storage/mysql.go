package storage

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/composition/orm"
	"gorm.io/gorm"
	"infoserver/lyric"
)

type mysqlStorageImpl struct {
	gorm orm.Gorm `aperture:""`
}

func MysqlImpl(gorm orm.Gorm) lyric.LyricStorageRepository {
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
	return lyric.Lens.ImplName("LyricStorageRepository", "mysql")
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
		result[i] = lyric.Song{
			Title: song.Title,
		}
		for _, artst := range song.Artists {
			result[i].Artist = append(result[i].Artist, artst.Name)
		}
		for _, lrc := range song.Lyrics {
			result[i].Lyrics = append(result[i].Lyrics, lyric.Lyric{Lang: lrc.Lang, Lyric: lrc.LyricText})
		}
	}

	return result, nil
}

func (m *mysqlStorageImpl) Search(keyword string) (result []lyric.Song, err error) {
	var songs []tableSong
	//SELECT DISTINCT song_id, title
	//FROM lyric_songs s
	//JOIN lyric_song_artists sa ON s.song_id = sa.table_song_song_id
	//JOIN lyric_artists a ON sa.table_song_song_id = a.artist_id
	//WHERE INSTR(LOWER('RΞOL · 染 Somari 【中日歌词】'), LOWER(s.title)) > 0 OR INSTR(LOWER('reol'), LOWER(a.name)) > 0;
	err = m.gorm.DB().
		Preload("Artists").
		Preload("Lyrics").
		Distinct("lyric_songs.song_id").
		Select("lyric_songs.song_id, lyric_songs.title").
		Joins("JOIN lyric_song_artists ON lyric_song_artists.table_song_song_id = lyric_songs.song_id").
		Joins("JOIN lyric_artists ON lyric_artists.artist_id = lyric_song_artists.table_artist_artist_id").
		Where("INSTR(LOWER(?), LOWER(lyric_songs.title)) > 0 OR INSTR(LOWER(?), LOWER(lyric_artists.name)) > 0", keyword, keyword).
		Limit(20).
		Find(&songs).Error

	if err != nil {
		return nil, err
	}

	result = make([]lyric.Song, len(songs))
	result = make([]lyric.Song, len(songs))
	for i, song := range songs {
		result[i] = lyric.Song{
			Title: song.Title,
		}
		for _, artst := range song.Artists {
			result[i].Artist = append(result[i].Artist, artst.Name)
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
			var existingSong tableSong
			// Load the song with its artists to check both name and associated artists
			result := tx.Preload("Artists").Where("title = ?", song.Title).First(&existingSong)

			if result.RowsAffected > 0 && artistsMatch(existingSong.Artists, song.Artist) {
				continue
			}

			// Either no song was found, or it has different artists
			newSong := tableSong{
				Title:   song.Title,
				Artists: []tableArtist{},
				Lyrics:  []tableLyric{},
			}

			for _, artistName := range song.Artist {
				var artist tableArtist
				// Check if artist exists, use existing artist if so
				if tx.Where("name = ?", artistName).First(&artist).RowsAffected == 0 {
					// Artist does not exist, create new
					artist = tableArtist{Name: artistName}
					tx.Create(&artist) // Create the artist
				}
				newSong.Artists = append(newSong.Artists, artist)
			}

			// Add lyrics as usual
			for _, lrc := range song.Lyrics {
				newSong.Lyrics = append(newSong.Lyrics, tableLyric{Lang: lrc.Lang, LyricText: lrc.Lyric})
			}

			// Save the new song
			if err := tx.Create(&newSong).Error; err != nil {
				return err // Rollback the transaction on error
			}
		}
		return nil
	})
	return err
}
