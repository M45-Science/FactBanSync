package main

import "time"

func mainLoop() {
	var LastFetchBans = time.Now()
	var LastWatch = time.Now()
	var LastRefresh = time.Now()

	//Loop, checking for new bans
	for serverRunning {
		time.Sleep(time.Second)

		if time.Since(LastFetchBans).Seconds() >= float64(serverConfig.FetchBansSeconds) {
			LastFetchBans = time.Now()

			fetchBanLists()
		}
		if time.Since(LastWatch).Seconds() >= float64(serverConfig.WatchFileSeconds) {
			LastWatch = time.Now()
			if serverConfig.FactorioBanFile != "" {
				watchBanFile()
			}
		}
		if time.Since(LastRefresh).Minutes() >= float64(serverConfig.RefreshListMinutes) {
			LastRefresh = time.Now()

			updateServerList()
		}
	}
}
