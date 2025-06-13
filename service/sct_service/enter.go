package sct_service

import (
	"bytes"
	"encoding/json"
	"flame_clouds/config"
	"flame_clouds/global"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
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

	// 如果有图片链接，添加到消息中
	//if notification.ImageURL != "" {
	//	message += fmt.Sprintf("  ![火烧云预测](%s)", notification.ImageURL)
	//}

	if !global.Config.Bot.Enable {
		logrus.Infof("未配置消息推送渠道 %s", message)
		return
	}
	logrus.Infof(message)

	type FtReq struct {
		Title string `json:"title"`
		Desp  string `json:"desp"`
	}
	byteData, _ := json.Marshal(FtReq{
		Title: notification.Event.EventType.String() + "火烧云预警",
		Desp:  message,
	})
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://sctapi.ftqq.com/%s.send", global.Config.Bot.SendKey), bytes.NewReader(byteData))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("消息推送失败 %s", err)
		return
	}
	byteData, _ = io.ReadAll(res.Body)

	type ServerType struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Pushid  string `json:"pushid"`
			Readkey string `json:"readkey"`
			Error   string `json:"error"`
			Errno   int    `json:"errno"`
		} `json:"data"`
	}

	var t ServerType
	err = json.Unmarshal(byteData, &t)
	if err != nil {
		logrus.Errorf("%s json解析失败, %s", string(byteData), err)
		return
	}

	logrus.Infof("推送成功 pushID:%s", t.Data.Pushid)
	return
}
