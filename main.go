package main

import (
	"flag"
	"log"
	"runtime/debug"
)

func main() {
	var runWizard bool

	//KB, MB
	debug.SetMemoryLimit(1024 * 1024 * 100)

	//Launch arguments
	flag.StringVar(&configPath, "configPath", defaultConfigPath, "config file path")
	var makeConfig bool
	flag.BoolVar(&makeConfig, "makeConfig", false, "make a default config file")
	flag.BoolVar(&runWizard, "runWizard", false, "run the setup wizard")
	var forceFetch bool
	flag.BoolVar(&forceFetch, "forceFetch", false, "force startup fetching ban lists from remotes")
	var verboseLogging bool
	flag.BoolVar(&verboseLogging, "verboseLogging", false, "force enable verbose logging")
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
	if verboseLogging {
		serverConfig.ServerPrefs.VerboseLogging = true
	}
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
	readServerBanList()
	compositeBans()
	updateWebCache()

	startWebserver()

	mainLoop()
}
