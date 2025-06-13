package flags

import "flag"

type FlagType struct {
	File string
}

var Options FlagType

func init() {
	flag.StringVar(&Options.File, "f", "settings.yaml", "配置文件的地址")
	flag.Parse()
}

func Run() {

}
