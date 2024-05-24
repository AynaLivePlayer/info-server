package delivery

import (
	sgin "github.com/rhine-tech/scene/scenes/gin"
	"infoserver/blivedm"
)

type appContext struct {
	dmSrv     blivedm.WebDanmuService      `aperture:""`
	openblive blivedm.OpenBLiveApiService  `aperture:""`
	connlog   blivedm.ConnectionLogService `aperture:""`
}

func GinApp(dmSrv blivedm.WebDanmuService,
	openblive blivedm.OpenBLiveApiService,
	connlog blivedm.ConnectionLogService,
) sgin.GinApplication {
	return &sgin.AppRoutes[appContext]{
		AppName:  blivedm.Lens.ImplNameNoVer("GinApplication"),
		BasePath: blivedm.Lens.String(),
		Actions: []sgin.Action[*appContext]{
			new(dmInfoV1Request),
			new(dmInfoRequest),
			new(openbliveAppStartRequest),
			new(openbliveAppEndRequest),
			new(openbliveHeartBeatRequest),
			new(connLogRequest),
			new(roomLogRequest),
		},
		Context: appContext{
			dmSrv:     dmSrv,
			openblive: openblive,
			connlog:   connlog,
		},
	}
}
