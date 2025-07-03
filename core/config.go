package core

import (
	"flame_clouds/config"
	"flame_clouds/config/types"
	"flame_clouds/flags"
	"flame_clouds/global"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

// ReadConfig 读取yaml配置
func ReadConfig() *config.Config {
	byteData, err := os.ReadFile(flags.Options.File)
	if err != nil {
		logrus.Fatalf("配置文件错误 %s", err)
	}
	var server config.Config
	err = yaml.Unmarshal(byteData, &server)
	if err != nil {
		logrus.Fatalf("配置文件格式错误 %s", err)
	}
	server.Monitor.Evening.EventType = types.Evening
	server.Monitor.Morning.EventType = types.Morning
	logrus.Infof("%s 配置文件加载成功", flags.Options.File)
	ConfigValid(&server)
	return &server
}

// ConfigValid 校验配置文件，因为要兼容老配置
func ConfigValid(c *config.Config) {
	logrus.Infof("验证配置文件")

	// 如果配置了老的city，那新的cityList就不生效
	if c.Monitor.City == "" {
		if len(c.Monitor.CityList) == 0 {
			logrus.Fatalf("未配置监控城市，程序退出")
		}
		logrus.Infof("监控城市 %v", c.Monitor.CityList)
	} else {
		if len(c.Monitor.CityList) != 0 {
			logrus.Warnf("存在老的city配置，新的cityList无效")
		}
		logrus.Infof("监控城市 %v", c.Monitor.City)
	}

	// 如果配置了老的消息推送，新的消息推送同样不生效
	if !c.Bot.Enable {
		logrus.Warnf("未配置消息通知渠道")
		return
	}
	if c.Bot.Target != "" {
		if len(c.Bot.TargetList) != 0 {
			logrus.Warnf("存在老的消息推送，新的消息推送无效")
		}
	}
}

func DumpConfig() {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		logrus.Errorf("转换错误 %s", err)
		return
	}

	err = os.WriteFile(flags.Options.File, byteData, 0666)
	if err != nil {
		logrus.Errorf("配置文件写入失败 %s", err)
		return
	}
	logrus.Infof("配置文件写入成功")

}
