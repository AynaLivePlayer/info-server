package delivery

import (
	"errors"
	"fmt"
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/infrastructure/asynctask"
	"github.com/rhine-tech/scene/infrastructure/datasource"
	"github.com/rhine-tech/scene/infrastructure/logger"
	"github.com/rhine-tech/scene/lens/permission"
	"scene-service/meowsic"
	"scene-service/onebot"
	"scene-service/onebot/protocol"
	"scene-service/onebot/scene/middleware"
	"strings"
)

type chatbotApp struct {
	redis      datasource.RedisDataSource `aperture:""`
	dispatcher asynctask.TaskDispatcher   `aperture:""`
	srv        meowsic.IMeowsicService    `aperture:""`
	log        logger.ILogger             `aperture:""`
}

func NewChatBotPlugin() onebot.ChatBotPlugin {
	return &chatbotApp{}
}

func (c *chatbotApp) Name() scene.ImplName {
	return meowsic.Lens.ImplNameNoVer("ChatBotPlugin")
}

func (c *chatbotApp) ImplName() scene.ImplName {
	return c.Name()
}

func (c *chatbotApp) Info() onebot.PluginInfo {
	return onebot.PluginInfo{
		Name:        "Meowsic",
		Description: "Meowsic",
	}
}

func (c *chatbotApp) Unload() error {
	return nil
}
func (c *chatbotApp) Load(router onebot.ChatBotEventRouter) error {
	router.HandlePrivateMessage(c.msgHandler, middleware.AuthContext())
	router.HandleGroupMessage(c.groupMsgHandler, middleware.AuthContext(), middleware.PermContext())
	return nil
}

func (c *chatbotApp) msgHandler(ctx scene.Context, client protocol.IOneBotApi, message *protocol.EventMessagePrivate) bool {
	return false
}

func (c *chatbotApp) groupMsgHandler(ctx scene.Context, client protocol.IOneBotApi, message *protocol.EventMessageGroup) bool {
	parts := strings.Fields(message.RawMessage)
	if len(parts) < 2 || parts[0] != "meowsic" {
		return false
	}

	switch parts[1] {
	case "providers":
		return c.handleProvidersCommand(ctx, client, message)
	case "status":
		if len(parts) < 3 {
			_, _ = client.SendGroupMessage(message.GroupID, []protocol.Message{
				protocol.NewMessageText("请指定音源提供者名称，例如: meowsic status netease"),
			}, false)
			return true
		}
		return c.handleStatusCommand(ctx, client, message, parts[2])
	case "updatestatus":
		if len(parts) < 3 {
			_, _ = client.SendGroupMessage(message.GroupID, []protocol.Message{
				protocol.NewMessageText("请指定音源提供者名称，例如: meowsic updatestatus netease"),
			}, false)
			return true
		}
		return c.handleUpdateStatusCommand(ctx, client, message, parts[2])
	}

	return false
}

func (c *chatbotApp) handleUpdateStatusCommand(ctx scene.Context, client protocol.IOneBotApi, message *protocol.EventMessageGroup, provider string) bool {
	// 调用服务更新状态
	err := c.srv.WithContext(ctx).UpdateStatus(provider)
	if err != nil {
		c.log.Errorf("更新音源状态失败: %v", err)
		if errors.Is(err, permission.ErrPermissionDenied) {
			_, _ = client.SendGroupMessage(message.GroupID, []protocol.Message{
				protocol.NewMessageText(fmt.Sprintf("无权限: %s", err.Error())),
			}, false)
			return true
		}
		_, _ = client.SendGroupMessage(message.GroupID, []protocol.Message{
			protocol.NewMessageText(fmt.Sprintf("更新音源状态失败: %s", err)),
		}, false)
		return true
	}
	_, _ = client.SendGroupMessage(message.GroupID, []protocol.Message{
		protocol.NewMessageText(fmt.Sprintf("✅ 已成功队列更新 %s 状态计划", provider)),
	}, false)
	return true
}

func (c *chatbotApp) handleProvidersCommand(ctx scene.Context, client protocol.IOneBotApi, message *protocol.EventMessageGroup) bool {
	providers, err := c.srv.WithContext(ctx).ListAllProvider()
	if err != nil {
		c.log.Errorf("获取音源提供者列表失败: %v", err)
		_, _ = client.SendGroupMessage(message.GroupID, []protocol.Message{
			protocol.NewMessageText("获取音源提供者列表失败"),
		}, false)
		return true
	}

	text := "可用音源提供者:\n" + strings.Join(providers, "\n")
	_, _ = client.SendGroupMessage(message.GroupID, []protocol.Message{
		protocol.NewMessageText(text),
	}, false)
	return true
}

func (c *chatbotApp) handleStatusCommand(ctx scene.Context, client protocol.IOneBotApi, message *protocol.EventMessageGroup, provider string) bool {
	status, err := c.srv.WithContext(ctx).GetStatus(provider)
	if err != nil {
		c.log.Errorf("获取音源状态失败: %v", err)
		_, _ = client.SendGroupMessage(message.GroupID, []protocol.Message{
			protocol.NewMessageText(fmt.Sprintf("获取音源状态失败: %s", err)),
		}, false)
		return true
	}

	// 创建状态符号
	statusSymbol := "✅"
	if !status.LoggedIn {
		statusSymbol = "⚠️"
	}

	updateTime := status.UpdateTime.Format("2006-01-02 15:04:05")

	// 创建API状态文本
	apiStatus := fmt.Sprintf(
		"%s API可用性 (登陆状态: %s)\n"+
			"    搜索API: %s %s\n"+
			"    信息API: %s %s\n"+
			"    文件API: %s %s\n"+
			"    VIP文件API: %s %s\n"+
			"    歌词API: %s %s\n"+
			"    播放列表API: %s %s\n"+
			"最后更新时间: %s",
		provider,
		statusSymbol,
		getStatusSymbol(status.Search.Status), getStatusText(status.Search),
		getStatusSymbol(status.Info.Status), getStatusText(status.Info),
		getStatusSymbol(status.FileUrl.Status), getStatusText(status.FileUrl),
		getStatusSymbol(status.VipFileUrl.Status), getStatusText(status.VipFileUrl),
		getStatusSymbol(status.Lyrics.Status), getStatusText(status.Lyrics),
		getStatusSymbol(status.Playlist.Status), getStatusText(status.Playlist),
		updateTime,
	)

	_, _ = client.SendGroupMessage(message.GroupID, []protocol.Message{
		protocol.NewMessageText(apiStatus),
	}, false)
	return true
}

// 辅助函数：获取状态符号
func getStatusSymbol(status bool) string {
	if status {
		return "✅"
	}
	return "❌"
}

// 辅助函数：获取状态文本
func getStatusText(status meowsic.ApiStatus) string {
	if status.Status {
		return "可用"
	}
	return "不可用"
}
