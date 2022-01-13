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

	//Read config file, rewrite for indentation/cleanup
	readConfigFile()
	writeConfigFile()

	//Logging
	startLog()
	log.Println(fmt.Sprintf("FactBanSync v%v", version))

	//Read our ban list, rewrite for indentation/cleanup
	readServerBanList()
	writeBanListFile()

	//Read our server list, then update it
	readServerListFile()
	updateServerList()

	//Run a webserver, if requested
	//TODO offer HTTPs with directions to make cert
	if serverConfig.RunWebServer {
		go func(WebPort int) {
			http.HandleFunc("/", handleFileRequest)
			http.ListenAndServe(":"+strconv.Itoa(serverConfig.WebPort), nil)
		}(serverConfig.WebPort)
		log.Println("Web server started http://localhost:" + strconv.Itoa(serverConfig.WebPort) + "/server-banlist.json")
	}

	var LastFetchBans = time.Now()
	var LastWatch = time.Now()
	var LastRefresh = time.Now()

	//Loop, checking for new bans
	for serverRunning {
		time.Sleep(time.Second)

		if time.Since(LastFetchBans).Seconds() >= float64(serverConfig.FetchBansInterval) {
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

func handleFileRequest(w http.ResponseWriter, r *http.Request) {
	defer time.Sleep(time.Millisecond * 100) //Max 10 requests per second

	if r.URL.Path == "/server-banlist.json" {
		if cachedBanList == "" {
			fmt.Fprintf(w, "No ban data")
			return
		}

		fmt.Fprintf(w, cachedBanList)

	} else {
		fmt.Fprintf(w, "404: File not found")
	}
}
