package service

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/infrastructure/asynctask"
	"github.com/rhine-tech/scene/infrastructure/logger"
	"infoserver/lyric"
)

type lyricService struct {
	storage    lyric.LyricStorageRepository
	providers  []lyric.LyricProvider
	logger     logger.ILogger           `aperture:""`
	dispatcher asynctask.TaskDispatcher `aperture:""`
}

func NewLyricService(storage lyric.LyricStorageRepository, providers ...lyric.LyricProvider) lyric.LyricService {
	return &lyricService{
		storage:   storage,
		providers: providers,
	}
}

func (l *lyricService) Setup() error {
	return nil
}

func (l *lyricService) SrvImplName() scene.ImplName {
	return lyric.Lens.ImplNameNoVer("LyricService")
}

func (l *lyricService) GetLyric(title string, artist string) (result []lyric.Song, err error) {
	l.logger.Infof("trying to searching lyric for %s - %s", title, artist)
	result, err = l.storage.GetLyric(title, artist)
	if err == nil && len(result) > 0 {
		return result, err
	}
	result = make([]lyric.Song, 0)
	for _, v := range l.providers {
		rs, err := v.GetLyric(title, artist)
		if err != nil {
			continue
		}
		l.logger.Infof("provider %s found lyric for %s - %s", v.ImplName(), title, artist)
		result = append(result, rs)
		l.dispatcher.Run(
			func() error {
				if e := l.storage.Add([]lyric.Song{rs}); e != nil {
					l.logger.Warnf("failed to add lyric %s (%s) to storage: %s", rs.Title, rs.Artist, e.Error())
				}
				return nil
			})
		break
	}
	if len(result) == 0 {
		l.logger.Warnf("failed to find lyric for %s - %s", title, artist)
		return result, lyric.ErrCantFindLyric.WithDetailStr(title + " - " + artist)
	}
	return result, nil
}

func (l *lyricService) Search(keyword string) (result []lyric.Song, err error) {
	l.logger.Infof("trying to search lyric for %s", keyword)
	result, err = l.storage.Search(keyword)
	if err != nil {
		l.logger.Warnf("failed to search lyric for %s: %s", keyword, err.Error())
		return nil, lyric.ErrSearchLyricFailed
	}
	return rankMedia(keyword, result), nil
}
