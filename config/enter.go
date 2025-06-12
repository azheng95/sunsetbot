package config

type Config struct {
	Hsy       Hsy       `yaml:"hsy"`
	ServerBot ServerBot `yaml:"serverBot"`
}

type Hsy struct {
	CheckAod float64 `yaml:"checkAod"`
	City     string  `yaml:"city"`
	WxDate   string  `yaml:"wxDate"`
	ZxDate   string  `yaml:"zxDate"`
}
type ServerBot struct {
	Enable  bool   `yaml:"enable"`
	SendKey string `yaml:"sendKey"`
}
