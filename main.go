package main

import (
	"flag"
	"log"
)

func main() {
	var runWizard bool

	//Launch arguments
	flag.StringVar(&configPath, "configPath", defaultConfigPath, "config file path")
	var makeConfig bool
	flag.BoolVar(&makeConfig, "makeConfig", false, "make a default config file")
	flag.BoolVar(&runWizard, "runWizard", false, "run the setup wizard")
	var forceFetch bool
	flag.BoolVar(&forceFetch, "forceFetch", false, "force startup fetching ban lists from remotes")
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
	log.Printf("FactBanSync v%v\n", ProgVersion)

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

	startWebserver()

	mainLoop()
}
