package service

import (
	"errors"
	"github.com/AynaLivePlayer/miaosic"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/composition/orm"
	"github.com/rhine-tech/scene/infrastructure/asynctask"
	"github.com/rhine-tech/scene/infrastructure/logger"
	"github.com/rhine-tech/scene/model/query"
	"scene-service/meowsic"
)

type meowsicService struct {
	updaters       map[string]meowsic.ISourceStatusUpdater
	sessionStorage meowsic.ISessionStorage                     `aperture:""`
	statusStorage  orm.GenericRepository[meowsic.SourceStatus] `aperture:""`
	log            logger.ILogger                              `aperture:""`
	cron           asynctask.CronTaskDispatcher                `aperture:""`
	taskDispatcher asynctask.TaskDispatcher                    `aperture:""`
	cronTask       asynctask.CronTask
}

func NewMeowsicService(updaters ...meowsic.ISourceStatusUpdater) meowsic.IMeowsicService {
	updatersMap := make(map[string]meowsic.ISourceStatusUpdater)
	for _, updater := range updaters {
		updatersMap[updater.ProviderName()] = updater
	}
	return &meowsicService{
		updaters: updatersMap,
		cronTask: asynctask.CronTask{
			Name: "meowsic.IMeowsicService.update_status",
		},
	}
}

func (m *meowsicService) Setup() error {
	m.cronTask.Func = func() error {
		m.updateAllStatus()
		return nil
	}
	m.recoverSessions()
	err := m.cron.AddTask("0 0 * * * ? ", &m.cronTask)
	if err != nil {
		return err
	}
	//m.updateAllStatus()
	return nil
}

func (m *meowsicService) SrvImplName() scene.ImplName {
	return meowsic.Lens.ImplNameNoVer("IMeowsicService")
}

func (m *meowsicService) ListAllProvider() ([]string, error) {
	return miaosic.ListAvailableProviders(), nil
}

func (m *meowsicService) GetStatus(provider string) (meowsic.SourceStatus, error) {
	first, b, err := m.statusStorage.FindFirst(query.Field("provider").Equal(provider), query.Field("update_time").Descending())
	if !b {
		m.log.WarnW("get status failed", "provider", provider, "err", err)
		return meowsic.SourceStatus{}, meowsic.ErrStatusNotFound.WithDetailStr(provider)
	}
	return first, nil
}

func (m *meowsicService) UpdateStatus(provider string) error {
	updater, ok := m.updaters[provider]
	if !ok {
		return meowsic.ErrSourceNotFound
	}
	m.taskDispatcher.Run(func() error {
		status, err := updater.Run()
		if err != nil {
			m.log.WarnW("update status failed", "provider", provider, "err", err)
			return nil
		}
		err = m.statusStorage.Create(&status)
		if err != nil {
			m.log.ErrorW("update status failed", "provider", provider, "err", err)
			return nil
		}
		updatedSession, err := updater.GetCurrentSession()
		if errors.Is(err, meowsic.ErrSourceNotLoginable) {
			return nil
		}
		if err != nil {
			m.log.WarnW("failed to get current session", "provider", updater.ProviderName(), "err", err)
			return nil
		}
		err = m.sessionStorage.UpdateSession(updater.ProviderName(), updatedSession)
		if err != nil {
			m.log.ErrorW("failed to update session", "provider", updater.ProviderName(), "err", err)
		}
		return nil
	})
	return nil
}

func (m *meowsicService) RestoreLogin(provider, session string) error {
	updater, ok := m.updaters[provider]
	if !ok {
		return meowsic.ErrSourceNotFound
	}
	err := updater.RestoreLogin(session)
	if err != nil {
		m.log.ErrorW("restore login failed", "provider", provider, "err", err)
		return meowsic.ErrorRestoreLoginFailed.WrapIfNot(err)
	}
	err = m.sessionStorage.UpdateSession(provider, session)
	if err != nil {
		m.log.ErrorW("restore login failed", "provider", provider, "err", err)
		return meowsic.ErrorRestoreLoginFailed
	}
	return nil
}

func (m *meowsicService) Logout(provider string) error {
	updater, ok := m.updaters[provider]
	if !ok {
		return meowsic.ErrSourceNotFound
	}
	err := updater.Logout()
	if err != nil {
		m.log.ErrorW("logout failed", "provider", provider, "err", err)
		return meowsic.ErrLogoutFailed
	}
	err = m.sessionStorage.RemoveSession(provider)
	if err != nil {
		m.log.ErrorW("RemoveSession failed", "provider", provider, "err", err)
		return meowsic.ErrLogoutFailed
	}
	return nil
}

func (m *meowsicService) recoverSessions() {
	for _, updater := range m.updaters {
		first, b2, err := m.sessionStorage.GetSession(updater.ProviderName())
		if err == nil && b2 {
			err = updater.RestoreLogin(first)
			if err != nil {
				m.log.ErrorW("failed to restore session", "provider", updater.ProviderName(), "err", err)
			} else {
				m.log.Info("restored session success", "provider", updater.ProviderName())
			}
		}
	}
}

func (m *meowsicService) updateAllStatus() {
	for _, updater := range m.updaters {
		m.log.InfoW("running status updater", "provider", updater.ProviderName())
		status, err := updater.Run()
		if err != nil {
			m.log.WarnW("failed to update status", "provider", updater.ProviderName(), "err", err)
			continue
		}
		err = m.statusStorage.Create(&status)
		if err != nil {
			m.log.ErrorW("failed to update status", "provider", updater.ProviderName(), "err", err)
			continue
		}
		updatedSession, err := updater.GetCurrentSession()
		if errors.Is(err, meowsic.ErrSourceNotLoginable) {
			continue
		}
		if err != nil {
			m.log.WarnW("failed to get current session", "provider", updater.ProviderName(), "err", err)
			continue
		}
		err = m.sessionStorage.UpdateSession(updater.ProviderName(), updatedSession)
		if err != nil {
			m.log.ErrorW("failed to update session", "provider", updater.ProviderName(), "err", err)
		}
	}
}

func (m *meowsicService) GetMediaUrl(providerName, identifier, quality string) ([]miaosic.MediaUrl, error) {
	provider, ok := miaosic.GetProvider(providerName)
	if !ok {
		return nil, meowsic.ErrSourceNotFound
	}
	_, ok = provider.MatchMedia(identifier)
	if !ok {
		return nil, meowsic.ErrInvalidMeta
	}
	return provider.GetMediaUrl(miaosic.MetaData{
		Provider:   providerName,
		Identifier: identifier,
	}, miaosic.Quality(quality))
}
