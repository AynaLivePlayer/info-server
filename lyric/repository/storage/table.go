package storage

import "infoserver/lyric"

type tableSong struct {
	SongID uint64 `json:"song_id" gorm:"unique;primaryKey;autoIncrement:true"`
	Title  string `gorm:"size:255"`

	// Associations
	Artists []tableArtist `gorm:"many2many:lyric_song_artists;"`
	Lyrics  []tableLyric  `gorm:"foreignKey:SongID"`
}

func (tableSong) TableName() string {
	return lyric.Lens.TableName("songs")
}

type tableArtist struct {
	ArtistID uint64 `json:"artist_id" gorm:"unique;primaryKey;autoIncrement:true"`
	Name     string `json:"name" gorm:"size:255"`

	// Associations
	Songs []tableSong `gorm:"many2many:lyric_song_artists;"`
}

func (tableArtist) TableName() string {
	return lyric.Lens.TableName("artists")
}

type tableLyric struct {
	LyricID   uint   `gorm:"primaryKey;autoIncrement"`
	SongID    uint64 `gorm:""`
	Lang      string `gorm:"size:50"`
	LyricText string `gorm:"type:text"`

	// Associations
	Song tableSong
}

func (tableLyric) TableName() string {
	return lyric.Lens.TableName("lyrics")
}

/*
SELECT
	DISTINCT s.song_id,
    s.Title AS Song_Name,
    a.Name AS Artist_Name,
    l.lyric_text AS Lyric,
    l.Lang AS Lyric_Language
FROM
    lyric_songs s
JOIN
    lyric_song_artists sa ON s.song_id = sa.table_song_song_id
JOIN
    lyric_artists a ON sa.table_artist_artist_id = a.artist_id
JOIN
    lyric_lyrics l ON s.song_id = l.song_id
WHERE
    s.Title = '染' AND (a.Name = '' OR a.Name = 'reol');
*/
