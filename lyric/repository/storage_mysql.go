package repository

import (
	"github.com/rhine-tech/scene/drivers/repos"
	"github.com/rhine-tech/scene/lens/infrastructure/datasource"
	"gorm.io/gorm"
	"infoserver/lyric"
)

type internalModel struct {
	Title  string `gorm:"column:title"`
	Artist string `gorm:"column:artist"`
	Lyric  string `gorm:"column:lyric"`
	Lang   string `gorm:"column:lang"`
}

func (internalModel) TableName() string {
	return "lyrics"
}

type mysqlStorageRepository struct {
	ds   datasource.MysqlDataSource `aperture:""`
	gorm *repos.GormRepo[internalModel]
}

func NewMysqlStorageRepository() lyric.LyricStorageRepository {
	return &mysqlStorageRepository{}
}

func (m *mysqlStorageRepository) Setup() (err error) {
	m.gorm, err = repos.UseGormMysql[internalModel](m.ds)
	return err
}

func (m *mysqlStorageRepository) RepoImplName() string {
	return "lyric.repository.LyricStorageRepository.mysql"
}

func (m *mysqlStorageRepository) Status() error {
	return m.ds.Status()
}

func (m *mysqlStorageRepository) Add(song lyric.Song) (err error) {
	models := make([]internalModel, 0)

	for _, v := range song.Lyrics {
		models = append(models, internalModel{
			Title:  song.Title,
			Artist: song.Artist,
			Lyric:  v.Lyric,
			Lang:   v.Lang,
		})
	}
	return m.gorm.DB.Create(&models).Error
}

func (m *mysqlStorageRepository) Search(name string, artist string) (result []lyric.Song, err error) {
	rs, err := m.gorm.FindPagination(func(db *gorm.DB) *gorm.DB {
		return db.Where("title = ?", name).Where("artist LIKE ?", "%"+artist+"%")
	}, 0, 20)
	if err != nil {
		return nil, err
	}
	tmpResult := make(map[string]*lyric.Song)
	for _, v := range rs.Results {
		key := v.Title + v.Artist
		if _, ok := tmpResult[key]; !ok {
			tmpResult[key] = &lyric.Song{
				Title:  v.Title,
				Artist: v.Artist,
				Lyrics: make([]lyric.Lyric, 0),
			}
		}
		tmpResult[key].Lyrics = append(tmpResult[key].Lyrics, lyric.Lyric{
			Lang:  v.Lang,
			Lyric: v.Lyric,
		})
	}
	result = make([]lyric.Song, 0)
	for _, v := range tmpResult {
		result = append(result, *v)
	}
	return result, nil
}
