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
		return
	}

	readConfigFile()
	writeConfigFile() //Clean up

	//Logging
	startLog()
	log.Println(fmt.Sprintf("FactBanSync v%v", ProgVersion))

	//Run a webserver, if requested
	exit := false
	if serverConfig.WebServer.RunWebServer {
		if serverConfig.PathData.FactorioBanFile == "" {
			log.Println("No factorio banlist file specified in config file")
			exit = true
		}
		if serverConfig.WebServer.SSLCertFile == "" || serverConfig.WebServer.SSLKeyFile == "" {
			log.Println("No SSL certificate or key file specified in config file")
			exit = true
		}
		if serverConfig.WebServer.DomainName == "" {
			log.Println("No domain name specified in config file")
			exit = true
		}
		if !exit {
			http.HandleFunc("/", handleFileRequest)
			server := &http.Server{
				Addr:         serverConfig.WebServer.DomainName + ":" + strconv.Itoa(serverConfig.WebServer.SSLWebPort),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 5 * time.Second,
				TLSConfig:    &tls.Config{ServerName: serverConfig.WebServer.DomainName},
			}
			go func(sc serverConfigData, serv *http.Server) {
				err := serv.ListenAndServeTLS(sc.WebServer.SSLCertFile, sc.WebServer.SSLKeyFile)
				if err != nil {
					log.Println(err)
				}
			}(serverConfig, server)
			if serverConfig.ServerPrefs.VerboseLogging {
				log.Println("Web server started:")
				log.Println(" https://" + serverConfig.WebServer.DomainName + ":" + strconv.Itoa(serverConfig.WebServer.SSLWebPort) + "/" + defaultFileWebName + ".gz")
				log.Println(" https://" + serverConfig.WebServer.DomainName + ":" + strconv.Itoa(serverConfig.WebServer.SSLWebPort) + "/" + defaultFileWebName)
			}
		} else {
			log.Println("Web server not started.")
		}
	}

	//Read banlist
	if serverConfig.PathData.FactorioBanFile != "" {
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
