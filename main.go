package main

import (
	"flame_clouds/core"
	"flame_clouds/global"
	"flame_clouds/service/cron_service"
)

func main() {
	global.Config = core.ReadConfig()
	core.InitLogger()

	cron_service.CronService()

	select {}
}
