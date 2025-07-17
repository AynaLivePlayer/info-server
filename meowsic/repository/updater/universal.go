package updater

import (
	"github.com/AynaLivePlayer/miaosic"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/infrastructure/logger"
	"scene-service/meowsic"
	"strings"
	"time"
)

type universalStatusUpdater struct {
	api           miaosic.MediaProvider
	provider      string
	hasVipApi     bool
	searchKeyword string
	infoMeta      string
	infoTitle     string
	fileMeta      string
	vipFileMeta   string
	lyricMeta     string
	playlistMeta  []string
	log           logger.ILogger `aperture:""`
}

func (b *universalStatusUpdater) Setup() error {
	if b.api == nil {
		b.api, _ = miaosic.GetProvider(b.provider)
	}

	return nil
}

func (b *universalStatusUpdater) ProviderName() string {
	return b.provider
}

func (b *universalStatusUpdater) ImplName() scene.ImplName {
	return meowsic.Lens.ImplName("ISourceStatusUpdater", b.provider)
}

func (b *universalStatusUpdater) RestoreLogin(session string) error {
	if b.api == nil {
		return meowsic.ErrSourceNotFound
	}
	loginableApi, ok := b.api.(miaosic.Loginable)
	if !ok {
		return meowsic.ErrSourceNotLoginable.WithDetailStr(b.provider)
	}
	return loginableApi.RestoreSession(session)
}

func (b *universalStatusUpdater) GetCurrentSession() (string, error) {
	if b.api == nil {
		return "", meowsic.ErrSourceNotFound
	}
	loginableApi, ok := b.api.(miaosic.Loginable)
	if !ok {
		return "", meowsic.ErrSourceNotLoginable.WithDetailStr(b.provider)
	}
	return loginableApi.SaveSession(), nil
}

func (b *universalStatusUpdater) Logout() error {
	if b.api == nil {
		return meowsic.ErrSourceNotFound
	}
	loginableApi, ok := b.api.(miaosic.Loginable)
	if !ok {
		return meowsic.ErrSourceNotLoginable.WithDetailStr(b.provider)
	}
	return loginableApi.Logout()
}

func (b *universalStatusUpdater) Run() (meowsic.SourceStatus, error) {
	if b.api == nil {
		var ok bool
		b.api, ok = miaosic.GetProvider(b.provider)
		if !ok {
			return meowsic.NewSourceStatusNotRegister(b.provider), nil
		}
	}
	// This branch will theoretically never enter, but just in case it's better to add the
	if b.api == nil {
		return meowsic.NewSourceStatusNotRegister(b.provider), nil
	}
	status := meowsic.SourceStatus{
		Provider:   b.provider,
		Search:     b.checkSearch(),
		Info:       b.checkInfo(),
		FileUrl:    b.checkFileUrl(),
		Lyrics:     b.checkLyrics(),
		Playlist:   b.checkPlaylist(),
		UpdateTime: time.Now(),
	}
	if b.hasVipApi {
		status.VipFileUrl = b.checkVipFileUrl()
	} else {
		status.VipFileUrl = status.FileUrl
	}
	loginableApi, ok := b.api.(miaosic.Loginable)
	if ok {
		status.LoggedIn = loginableApi.IsLogin()
	}
	return status, nil
}

func (b *universalStatusUpdater) checkSearch() meowsic.ApiStatus {
	result, err := b.api.Search(b.searchKeyword, 1, 20)
	if err != nil {
		return meowsic.NewFailApiStatus(err.Error())
	}
	if len(result) == 0 {
		return meowsic.NewFailApiStatus("search result is empty")
	}
	return meowsic.NewOkApiStatus()
}

func (b *universalStatusUpdater) checkInfo() meowsic.ApiStatus {
	meta, ok := b.api.MatchMedia(b.infoMeta)
	if !ok {
		return meowsic.NewFailApiStatus("source meta invalid, should not appear")
	}
	media, err := b.api.GetMediaInfo(meta)
	if err != nil {
		return meowsic.NewFailApiStatus(err.Error())
	}

	if media.Title != b.infoTitle || media.Title == "" {
		return meowsic.NewFailApiStatus("info not match, should not appear")
	}
	return meowsic.NewOkApiStatus()
}

func (b *universalStatusUpdater) checkFileUrl() meowsic.ApiStatus {
	meta, ok := b.api.MatchMedia(b.fileMeta)
	if !ok {
		return meowsic.NewFailApiStatus("source meta invalid, should not appear")
	}
	url, err := b.api.GetMediaUrl(meta, miaosic.QualityAny)
	if err != nil {
		return meowsic.NewFailApiStatus(err.Error())
	}
	if len(url) == 0 {
		return meowsic.NewFailApiStatus("no url was found")
	}
	if !strings.HasPrefix(url[0].Url, "http") {
		return meowsic.NewFailApiStatus("invalid url, should not appear")
	}
	return meowsic.NewOkApiStatus()
}

func (b *universalStatusUpdater) checkVipFileUrl() meowsic.ApiStatus {
	meta, ok := b.api.MatchMedia(b.vipFileMeta)
	if !ok {
		return meowsic.NewFailApiStatus("source meta invalid, should not appear")
	}
	url, err := b.api.GetMediaUrl(meta, miaosic.QualityAny)
	if err != nil {
		return meowsic.NewFailApiStatus(err.Error())
	}
	if len(url) == 0 {
		return meowsic.NewFailApiStatus("no url was found")
	}
	if !strings.HasPrefix(url[0].Url, "http") {
		return meowsic.NewFailApiStatus("invalid url, should not appear")
	}
	return meowsic.NewOkApiStatus()
}

func (b *universalStatusUpdater) checkLyrics() meowsic.ApiStatus {
	meta, ok := b.api.MatchMedia(b.vipFileMeta)
	if !ok {
		return meowsic.NewFailApiStatus("source meta invalid, should not appear")
	}
	result, err := b.api.GetMediaLyric(meta)
	if err != nil {
		return meowsic.NewFailApiStatus(err.Error())
	}
	if len(result) == 0 {
		return meowsic.NewFailApiStatus("no lyrics was found")
	}
	return meowsic.NewOkApiStatus()
}

func (b *universalStatusUpdater) checkPlaylist() meowsic.ApiStatus {
	for _, uri := range b.playlistMeta {
		meta, ok := b.api.MatchPlaylist(uri)
		if !ok {
			return meowsic.NewFailApiStatus("source meta invalid, api might not implement")
		}
		playlist, err := b.api.GetPlaylist(meta)
		if err != nil {
			return meowsic.NewFailApiStatus(err.Error())
		}
		if playlist.Title == "" {
			return meowsic.NewFailApiStatus("playlist title is empty")
		}
		if len(playlist.Medias) == 0 {
			return meowsic.NewFailApiStatus("playlist medias is empty")
		}
	}
	return meowsic.NewOkApiStatus()
}
