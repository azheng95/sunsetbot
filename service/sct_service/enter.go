package sct_service

import (
	"bytes"
	"encoding/json"
	"flame_clouds/global"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type AlertNotification struct {
	City      string  `json:"city"`
	AOD       float64 `json:"aod"`
	EventType string  `json:"event_type"`
	EventTime string  `json:"event_time"`
	ImageURL  string  `json:"image_url,omitempty"`
}

// SendServerNotification 发送server通知
func SendServerNotification(notification AlertNotification) error {
	// 构建消息内容
	var eventName string
	switch notification.EventType {
	case "set_1":
		eventName = "今日日落"
	case "rise_2":
		eventName = "明日日出"
	default:
		eventName = notification.EventType
	}

	message := fmt.Sprintf(
		"【火烧云预警】城市: %s  事件: %s  时间: %s  AOD指标: %.2f 满足拍摄条件!",
		notification.City, eventName, notification.EventTime, notification.AOD,
	)
	strings.ReplaceAll(message, "<br>", "")

	// 如果有图片链接，添加到消息中
	//if notification.ImageURL != "" {
	//	message += fmt.Sprintf("  ![火烧云预测](%s)", notification.ImageURL)
	//}

	if !global.Config.ServerBot.Enable {
		logrus.Infof("未配置消息推送渠道 %s", message)
		return nil
	}
	logrus.Infof(message)

	type FtReq struct {
		Title string `json:"title"`
		Desp  string `json:"desp"`
	}
	byteData, _ := json.Marshal(FtReq{
		Title: eventName + "火烧云预警",
		Desp:  message,
	})
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://sctapi.ftqq.com/%s.send", global.Config.ServerBot.SendKey), bytes.NewReader(byteData))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("消息推送失败 %s", err)
		return err
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
		return err
	}

	logrus.Infof("推送成功 pushID:%s", t.Data.Pushid)

	return nil
}
