package cron_service

import (
	"flame_clouds/global"
	"flame_clouds/service/hsy_service"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"time"
)

func CronService() {
	timezone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logrus.Infof("parse location Asia/Shanghai error, %s", err)
		timezone, err = time.LoadLocation("UTC+8")
		if err != nil {
			logrus.Infof("parse location UTC+8 error, %s", err)
			return
		}
	}

	// 看看配置文件是否正确
	if global.Config.Hsy.City == "" {
		logrus.Fatalf("未配置火烧云监控城市")
	}
	if global.Config.Hsy.CheckAod == 0 {
		logrus.Fatalf("未配置火烧云监控指标")
	}
	if !global.Config.ServerBot.Enable {
		logrus.Warnf("未配置消息通知渠道")
	}

	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	// 每天4点获取当天晚霞数据
	_, err = crontab.AddFunc(global.Config.Hsy.WxDate, func() {
		hsy_service.GetCitySunsetData("set_1")
	})
	if err != nil {
		logrus.Fatalf("晚霞时间配置错误 %s", err)
	}
	// 每天8点获取每天朝霞数据
	_, err = crontab.AddFunc(global.Config.Hsy.ZxDate, func() {
		hsy_service.GetCitySunsetData("rise_2")
	})
	if err != nil {
		logrus.Fatalf("朝霞时间配置错误 %s", err)
	}
	logrus.Infof("晚霞监控时间 %s", global.Config.Hsy.WxDate)
	logrus.Infof("朝霞监控时间 %s", global.Config.Hsy.ZxDate)
	crontab.Start()
}
