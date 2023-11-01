package service

import (
	"github.com/rhine-tech/scene/lens/infrastructure/asynctask"
	"github.com/rhine-tech/scene/lens/infrastructure/logger"
	"infoserver/lyric"
)

type lyricService struct {
	storage    lyric.LyricStorageRepository
	providers  []lyric.LyricSearchRepository
	logger     logger.ILogger           `aperture:""`
	dispatcher asynctask.TaskDispatcher `aperture:""`
}

func NewLyricService(storage lyric.LyricStorageRepository, providers ...lyric.LyricSearchRepository) lyric.LyricService {
	return &lyricService{
		storage:   storage,
		providers: providers,
	}
}

func (l *lyricService) Setup() error {
	l.logger = l.logger.WithPrefix(l.SrvImplName())
	return nil
}

func (l *lyricService) SrvImplName() string {
	return "lyric.service.LyricService"
}

func (l *lyricService) Search(title string, artist string) (result []lyric.Song, err error) {
	l.logger.Infof("trying to searching lyric for %s - %s", title, artist)
	result, err = l.storage.Search(title, artist)
	if err == nil && len(result) > 0 {
		return result, err
	}
	result = make([]lyric.Song, 0)
	for _, v := range l.providers {
		rs, err := v.Search(title, artist)
		if err != nil {
			continue
		}
		l.logger.Infof("provider %s found lyric for %s - %s", v.RepoImplName(), title, artist)
		result = append(result, rs)
		l.dispatcher.Run(
			func() error {
				if e := l.storage.Add(rs); e != nil {
					l.logger.Warnf("failed to add lyric %s to storage: %s", rs, e.Error())
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
