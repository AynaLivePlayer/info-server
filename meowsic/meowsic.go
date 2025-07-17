package meowsic

import (
	"github.com/AynaLivePlayer/miaosic"
	"github.com/rhine-tech/scene"
	"time"
)

const Lens scene.ModuleName = "meowsic"

type ApiStatus struct {
	Status     bool   `json:"status"`
	FailReason string `json:"fail_reason"`
}

func NewOkApiStatus() ApiStatus {
	return ApiStatus{
		Status: true,
	}
}

func NewFailApiStatus(reason string) ApiStatus {
	return ApiStatus{
		Status:     false,
		FailReason: reason,
	}
}

func newApiStatusSourceNotRegister() ApiStatus {
	return ApiStatus{
		Status:     false,
		FailReason: "source not register",
	}
}

type SourceStatus struct {
	Provider   string    `json:"provider"`
	LoggedIn   bool      `json:"logged_in"`
	Search     ApiStatus `json:"search" gorm:"embedded;embeddedPrefix:search_"`
	Info       ApiStatus `json:"info" gorm:"embedded;embeddedPrefix:info_"`
	FileUrl    ApiStatus `json:"file_url" gorm:"embedded;embeddedPrefix:file_url_"`
	VipFileUrl ApiStatus `json:"vip_file_url" gorm:"embedded;embeddedPrefix:vip_file_url_"`
	Lyrics     ApiStatus `json:"lyrics" gorm:"embedded;embeddedPrefix:lyrics_"`
	Playlist   ApiStatus `json:"playlist" gorm:"embedded;embeddedPrefix:playlist_"`
	UpdateTime time.Time `json:"update_time"`
}

func (d *SourceStatus) TableName() string {
	return Lens.TableName("status")
}

func NewSourceStatusNotRegister(provider string) SourceStatus {
	return SourceStatus{
		Provider:   provider,
		Search:     newApiStatusSourceNotRegister(),
		Info:       newApiStatusSourceNotRegister(),
		FileUrl:    newApiStatusSourceNotRegister(),
		VipFileUrl: newApiStatusSourceNotRegister(),
		Lyrics:     newApiStatusSourceNotRegister(),
		Playlist:   newApiStatusSourceNotRegister(),
		UpdateTime: time.Now(),
	}
}

type ISourceStatusUpdater interface {
	scene.Named
	ProviderName() string
	Run() (SourceStatus, error)
	RestoreLogin(session string) error
	GetCurrentSession() (string, error)
	Logout() error
}

type SourceSession struct {
	Provider string `json:"provider"`
	Session  string `json:"session"`
}

func (d *SourceSession) TableName() string {
	return Lens.TableName("sesssions")
}

type ISessionStorage interface {
	GetSession(provider string) (string, bool, error)
	RemoveSession(provider string) error
	UpdateSession(provider string, session string) error
}

type IMeowsicService interface {
	scene.Service
	scene.WithContext[IMeowsicService]
	ListAllProvider() ([]string, error)
	GetStatus(provider string) (SourceStatus, error)
	UpdateStatus(provider string) error
	RestoreLogin(provider, session string) error
	Logout(provider string) error

	// following are miaosic related api

	//SearchByProvider(provider string, keyword string, page, size int) ([]miaosic.MediaInfo, error)
	GetMediaUrl(provider, identifier, quality string) ([]miaosic.MediaUrl, error)
	//GetMediaInfo(provider, identifier string) (miaosic.MediaInfo, error)
	//GetMediaLyric(provider, identifier string) ([]miaosic.Lyrics, error)
}
