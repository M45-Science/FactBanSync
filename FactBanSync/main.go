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
	var runWizard bool

	//Launch arguments
	flag.StringVar(&configPath, "configPath", defaultConfigPath, "config file path")
	var makeConfig bool
	flag.BoolVar(&makeConfig, "makeConfig", false, "make a default config file")
	flag.BoolVar(&runWizard, "runWizard", false, "run the setup wizard")
	flag.Parse()

	//Make config file if requested
	if makeConfig {
		makeDefaultConfigFile()
		return
	}

	if runWizard {
		setupWizard()
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
		log.Println(" http://localhost:" + strconv.Itoa(serverConfig.WebPort) + "/" + defaultFileWebName + ".gz")
		log.Println(" http://localhost:" + strconv.Itoa(serverConfig.WebPort) + "/" + defaultFileWebName)
	}

	if serverConfig.FactorioBanFile != "" {
		readServerBanList()
	}
	readServerListFile()

	//Fetch if we don't have anything
	if len(serverList.ServerList) == 0 {
		updateServerList()
		fetchBanLists()
		compositeBans()
	}
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

//Web server
func handleFileRequest(w http.ResponseWriter, r *http.Request) {
	defer time.Sleep(time.Millisecond * 100) //Max 10 requests per second

	//Cached gzip copy
	if r.URL.Path == "/"+defaultFileWebName+".gz" {
		if cachedBanListGz == nil {
			noDataReply(w)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "gzip")
		cachedBanListLock.Lock()
		w.Write(cachedBanListGz)
		cachedBanListLock.Unlock()

		//Cached copy
	} else if r.URL.Path == "/"+defaultFileWebName {
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
