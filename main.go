package main

import (
	"flame_clouds/core"
	"flame_clouds/flags"
	"flame_clouds/global"
	"flame_clouds/service/cron_service"
)

func main() {
	global.Config = core.ReadConfig()
	core.InitLogger()

	flags.Run()

	cron_service.CronService()

	select {}
}
