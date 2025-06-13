package message_push_service

import (
	serverchan_sdk "github.com/easychen/serverchan-sdk-golang"
	"github.com/sirupsen/logrus"
)

type FtMsg struct {
	Key string
}

func (f FtMsg) Push(title string, des string) (error error) {
	resp, err := serverchan_sdk.ScSend(f.Key, title, des, nil)
	if err != nil {
		return
	}
	logrus.Infof("消息推送成功 %#v", resp)
	return
}
