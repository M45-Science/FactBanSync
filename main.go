package main

import (
	"flag"
	"fmt"
	"log"
)

const version = "0.0.1"

func main() {

	//Launch arguments
	configPath = *flag.String("configPath", defaultConfigPath, "config file path")
	flag.Parse()

	readConfigFile()
	writeConfigFile()

	startLog()
	log.Println(fmt.Sprintf("FactBanSync v%v", version))

	updateServerList()

	readServerListFile()
	readServerBanList()

	writeBanListFile() //To clean up formatting
}
