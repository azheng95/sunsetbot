package config

import "flame_clouds/config/types"

type Config struct {
	Monitor Monitor `yaml:"monitor"`
	Bot     Bot     `yaml:"bot"`
}

type Monitor struct {
	City    string       `yaml:"city"`
	Evening MonitorEvent `yaml:"evening"` // 晚霞
	Morning MonitorEvent `yaml:"morning"` // 朝霞
}

type MonitorEvent struct {
	EventType types.EventType `yaml:"-"`
	Enable    bool            `yaml:"enable"`
	CheckAod  float64         `yaml:"checkAod"`
	Time      string          `yaml:"time"`
}

type BotTargetType string

const FtBot BotTargetType = "ft"

type Bot struct {
	Enable  bool          `yaml:"enable"`
	Target  BotTargetType `yaml:"target"`
	SendKey string        `yaml:"sendKey"`
}
