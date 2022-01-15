package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const version = "0.0.1"

func main() {

	//Launch arguments
	flag.StringVar(&configPath, "configPath", defaultConfigPath, "config file path")
	var makeConfig bool
	flag.BoolVar(&makeConfig, "makeConfig", false, "make a default config file")
	flag.Parse()

	//Make config file if requested
	if makeConfig {
		makeDefaultConfigFile()
		return
	}

	readConfigFile()

	//Logging
	startLog()
	log.Println(fmt.Sprintf("FactBanSync v%v", version))

	//Run a webserver, if requested
	//TODO offer HTTPs with directions to make cert
	if serverConfig.RunWebServer {
		go func(WebPort int) {
			http.HandleFunc("/", handleFileRequest)
			http.ListenAndServe(":"+strconv.Itoa(serverConfig.WebPort), nil)
		}(serverConfig.WebPort)
		log.Println("Web server started:")
		log.Println(" http://localhost:" + strconv.Itoa(serverConfig.WebPort) + "/" + banFileWebName + ".gz")
		log.Println(" http://localhost:" + strconv.Itoa(serverConfig.WebPort) + "/" + banFileWebName)
	}

	readServerBanList()
	readServerListFile()

	//Startup update
	updateServerList()
	fetchBanLists()
	updateWebCache()

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
			if serverConfig.BanFile != "" {
				watchBanFile()
			}
		}
		if time.Since(LastRefresh).Minutes() >= float64(serverConfig.RefreshListMinutes) {
			LastRefresh = time.Now()

			updateServerList()
		}
	}
}

//Web server
func handleFileRequest(w http.ResponseWriter, r *http.Request) {
	defer time.Sleep(time.Millisecond * 100) //Max 10 requests per second

	//Cached gzip copy
	if r.URL.Path == "/"+banFileWebName+".gz" {
		if cachedBanListGz == nil {
			noDataReply(w)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "gzip")
		cachedBanListLock.Lock()
		w.Write(cachedBanListGz)
		cachedBanListLock.Unlock()

		//Cached copy
	} else if r.URL.Path == "/"+banFileWebName {
		if cachedBanList == nil {
			noDataReply(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		cachedBanListLock.Lock()
		w.Write(cachedBanList)
		cachedBanListLock.Unlock()
	} else {
		//Not found
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
		fmt.Fprintf(w, "404: File not found")
	}
}

func noDataReply(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	fmt.Fprintf(w, "No ban data")
}
