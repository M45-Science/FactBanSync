package main

import (
	"fmt"
	"strconv"
)

func setupWizard() {
	makeDefaultConfigFile()

	fmt.Println("Community or server name: ")

	var communityName string
	fmt.Scanln(&communityName)
	if communityName == "" {
		communityName = defaultName
	}
	serverConfig.Name = communityName

	fmt.Println("Run a webserver to provide server-banlist.json? (y/N)")

	var runWebServer string
	fmt.Scanln(&runWebServer)
	if runWebServer == "Y" || runWebServer == "y" {
		serverConfig.RunWebServer = true

		fmt.Println("Web server port (8080): ")

		var webPort string
		fmt.Scanln(&webPort)
		if webPort == "" {
			serverConfig.WebPort = defaultWebPort
		} else {
			serverConfig.WebPort, _ = strconv.Atoi(webPort)
		}
	} else {
		serverConfig.RunWebServer = false
	}

	fmt.Println("Auto-subscribe to new servers? (Y/n)")

	var autoSubscribe string
	fmt.Scanln(&autoSubscribe)

	if autoSubscribe == "N" || autoSubscribe == "n" {
		serverConfig.AutoSubscribe = false
	} else {
		serverConfig.AutoSubscribe = true
	}

	fmt.Println("Require reason for bans? (y/N)")

	var requireReason string
	fmt.Scanln(&requireReason)

	if requireReason == "Y" || requireReason == "y" {
		serverConfig.RequireReason = true
	} else {
		serverConfig.RequireReason = false
	}

	fmt.Println("Strip ban reasons from public ban list? (y/N)")

	var stripReasons string
	fmt.Scanln(&stripReasons)

	if stripReasons == "Y" || stripReasons == "y" {
		serverConfig.StripReasons = true
	} else {
		serverConfig.StripReasons = false
	}

	writeConfigFile()

}
