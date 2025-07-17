package service

import (
	"github.com/AynaLivePlayer/miaosic"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/lens/permission"
	"scene-service/meowsic"
)

type srvWithCtx struct {
	srv meowsic.IMeowsicService `aperture:""`
	ctx scene.Context
}

func (s *srvWithCtx) SrvImplName() scene.ImplName {
	return meowsic.Lens.ImplName("IMeowsicService", "withcontext")
}

func (s *srvWithCtx) ListAllProvider() ([]string, error) {
	return s.srv.ListAllProvider()
}

func (s *srvWithCtx) GetStatus(provider string) (meowsic.SourceStatus, error) {
	return s.srv.GetStatus(provider)
}

func (s *srvWithCtx) UpdateStatus(provider string) error {
	ok := permission.HasPermissionInCtx(s.ctx, meowsic.PermAdmin)
	if !ok {
		return permission.ErrPermissionDenied.WithDetailStr(meowsic.PermAdmin.String())
	}
	return s.srv.UpdateStatus(provider)
}

func (s *srvWithCtx) RestoreLogin(provider, session string) error {
	ok := permission.HasPermissionInCtx(s.ctx, meowsic.PermAdmin)
	if !ok {
		return permission.ErrPermissionDenied.WithDetailStr(meowsic.PermAdmin.String())
	}
	return s.srv.RestoreLogin(provider, session)
}

func (s *srvWithCtx) Logout(provider string) error {
	ok := permission.HasPermissionInCtx(s.ctx, meowsic.PermAdmin)
	if !ok {
		return permission.ErrPermissionDenied.WithDetailStr(meowsic.PermAdmin.String())
	}
	return s.srv.Logout(provider)
}

func (s *srvWithCtx) GetMediaUrl(provider, identifier, quality string) ([]miaosic.MediaUrl, error) {
	return s.srv.GetMediaUrl(provider, identifier, quality)
}

func (s *srvWithCtx) WithContext(ctx scene.Context) meowsic.IMeowsicService {
	return &srvWithCtx{srv: s.srv, ctx: ctx}
}

func (m *meowsicService) WithContext(ctx scene.Context) meowsic.IMeowsicService {
	return &srvWithCtx{srv: m, ctx: ctx}
}
