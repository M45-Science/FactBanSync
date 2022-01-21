package main

import (
	"crypto/tls"
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
	var forceFetch bool
	flag.BoolVar(&forceFetch, "forceFetch", false, "force fetching ban lists from remotes")
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
	writeConfigFile() //Clean up

	//Logging
	startLog()
	log.Println(fmt.Sprintf("FactBanSync v%v", version))

	//Run a webserver, if requested
	//TODO offer HTTPs with directions to make cert
	if serverConfig.RunWebServer {
		http.HandleFunc("/", handleFileRequest)
		server := &http.Server{
			Addr:         serverConfig.DomainName + ":" + strconv.Itoa(serverConfig.SSLWebPort),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			TLSConfig:    &tls.Config{ServerName: serverConfig.DomainName},
		}
		go func(sc serverConfigData, serv *http.Server) {
			serv.ListenAndServeTLS(sc.SSLCertFile, sc.SSLKeyFile)
		}(serverConfig, server)
		log.Println("Web server started:")
		log.Println(" https://" + serverConfig.DomainName + ":" + strconv.Itoa(serverConfig.SSLWebPort) + "/" + defaultFileWebName + ".gz")
		log.Println(" https://" + serverConfig.DomainName + ":" + strconv.Itoa(serverConfig.SSLWebPort) + "/" + defaultFileWebName)
	}

	if serverConfig.FactorioBanFile != "" {
		readServerBanList()
	}
	readServerListFile()
	readBanCache()

	//Fetch if we don't have anything
	if len(serverList.ServerList) == 0 || forceFetch {
		updateServerList()
		fetchBanLists()
	}
	compositeBans()
	updateWebCache()

	mainLoop()
}
