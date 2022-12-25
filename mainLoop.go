package main

import "time"

func mainLoop() {
	var LastFetchBans = time.Now()
	var LastWatch = time.Now()
	var LastRefresh = time.Now()

	//Loop, checking for new bans
	for serverRunning {
		time.Sleep(time.Second)

		if time.Since(LastFetchBans).Minutes() >= float64(serverConfig.ServerPrefs.FetchBansMinutes) {
			LastFetchBans = time.Now()

			fetchBanLists()
			compositeBans()
			updateWebCache()
		}
		//Don't run if no file specificed
		if serverConfig.PathData.FactorioBanFile != "" {
			if time.Since(LastWatch).Seconds() >= float64(serverConfig.ServerPrefs.WatchFileSeconds) {
				LastWatch = time.Now()
				if serverConfig.PathData.FactorioBanFile != "" {
					watchBanFile()
				}
			}
		}
		if time.Since(LastRefresh).Hours() >= float64(serverConfig.ServerPrefs.RefreshListHours) {
			LastRefresh = time.Now()

			updateServerList()

		}
	}
}
