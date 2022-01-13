package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

const version = "0.0.1"

func main() {

	//Launch arguments
	configPath = *flag.String("configPath", defaultConfigPath, "config file path")
	flag.Parse()

	readConfigFile()
	writeConfigFile() //To clean up formatting

	startLog()
	log.Println(fmt.Sprintf("FactBanSync v%v", version))

	readServerBanList()
	writeBanListFile() //To clean up formatting

	readServerListFile()
	updateServerList()

	var LastFetchBans = time.Now()
	var LastWatch = time.Now()
	var LastRefresh = time.Now()

	for serverRunning {
		time.Sleep(time.Second)

		if time.Since(LastFetchBans).Minutes() >= float64(serverConfig.FetchBansInterval) {
			LastFetchBans = time.Now()

			//Fetch bans (TODO)
		}
		if time.Since(LastWatch).Seconds() >= float64(serverConfig.WatchInterval) {
			LastWatch = time.Now()

			WatchBanFile()
		}
		if time.Since(LastRefresh).Minutes() >= float64(serverConfig.RefreshListInterval) {
			LastRefresh = time.Now()

			updateServerList()
		}
	}
}
