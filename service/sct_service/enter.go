package sct_service

import (
	"flame_clouds/config"
	"flame_clouds/global"
	"fmt"
	serverchan_sdk "github.com/easychen/serverchan-sdk-golang"
	"github.com/sirupsen/logrus"
	"strings"
)

type AlertNotification struct {
	City      string              `json:"city"`
	Quality   float64             `json:"quality"`
	Event     config.MonitorEvent `json:"event"`
	EventTime string              `json:"event_time"`
	ImageURL  string              `json:"image_url,omitempty"`
}

// SendServerNotification 发送server通知
func SendServerNotification(notification AlertNotification) {
	// 构建消息内容
	message := fmt.Sprintf(
		"【火烧云预警】城市: %s  事件: %s  时间: %s  火烧云质量: %.2f 满足拍摄条件!",
		notification.City, notification.Event.EventType.String(), notification.EventTime, notification.Quality,
	)
	message = strings.ReplaceAll(message, "<br>", "")

	if !global.Config.Bot.Enable {
		logrus.Infof("未配置消息推送渠道 %s", message)
		return
	}
	logrus.Infof(message)

	title := notification.Event.EventType.String() + "火烧云预警"

	resp, err := serverchan_sdk.ScSend(global.Config.Bot.SendKey, title, message, nil)
	if err != nil {
		logrus.Errorf("消息推送失败 %s", err)
		return
	}
	logrus.Infof("消息推送成功 %#v", resp)
	return
}
