package main

import (
	"fmt"
	"os"
	"strconv"
)

func setupWizard() {
	makeDefaultConfigFile()

	fmt.Println("Pressing enter on a question will accept the default value.")
	fmt.Println("Community or server name: (Used to skip ourselves in the server list)")

	var communityName string
	fmt.Scanln(&communityName)
	if communityName == "" {
		communityName = defaultCommunityName
	}
	serverConfig.CommunityName = communityName

	fmt.Println("You will need a certificate and key file for HTTPS.")
	fmt.Println("Run a rate-limited, cached HTTPS (SSL) webserver, to provide server-banlist.json? (y/N)")

	var runSSLWebServer string
	fmt.Scanln(&runSSLWebServer)
	if runSSLWebServer == "y" || runSSLWebServer == "Y" {
		serverConfig.WebServer.RunWebServer = true

		fmt.Println("HTTPS Web server port (8443): ")

		var webPort string
		fmt.Scanln(&webPort)
		if webPort == "" {
			serverConfig.WebServer.SSLWebPort = defaultSSLWebPort
		} else {
			serverConfig.WebServer.SSLWebPort, _ = strconv.Atoi(webPort)
		}

		fmt.Println("Domain name (REQUIRED): ")

		var domainName string
		fmt.Scanln(&domainName)
		if domainName == "" {
			domainName = "test.com"
		}
		serverConfig.WebServer.DomainName = domainName

	} else {
		serverConfig.WebServer.SSLWebPort = 0
		serverConfig.WebServer.RunWebServer = false
	}

	fmt.Println("Auto-subscribe to new servers? (Y/n)")

	var autoSubscribe string
	fmt.Scanln(&autoSubscribe)

	if autoSubscribe == "N" || autoSubscribe == "n" {
		serverConfig.ServerPrefs.AutoSubscribe = false
	} else {
		serverConfig.ServerPrefs.AutoSubscribe = true
	}

	fmt.Println("Require reason for bans? (y/N)")

	var requireReason string
	fmt.Scanln(&requireReason)

	if requireReason == "Y" || requireReason == "y" {
		serverConfig.ServerPrefs.RequireReason = true
	} else {
		serverConfig.ServerPrefs.RequireReason = false
	}

	fmt.Println("Strip ban reasons from public ban list? (y/N)")

	var stripReasons string
	fmt.Scanln(&stripReasons)

	if stripReasons == "Y" || stripReasons == "y" {
		serverConfig.ServerPrefs.StripReasons = true
	} else {
		serverConfig.ServerPrefs.StripReasons = false
	}

	fmt.Println("How often do you want to refresh the list of servers (in hours)? (12)")

	var refreshListHours string
	fmt.Scanln(&refreshListHours)
	if refreshListHours == "" {
		serverConfig.ServerPrefs.RefreshListHours = defaultRefreshListHours
	} else {
		serverConfig.ServerPrefs.RefreshListHours, _ = strconv.Atoi(refreshListHours)
	}

	fmt.Println("How often do you want to fetch ban lists from server (in minutes)? (15)")

	var fetchBansMinutes string
	fmt.Scanln(&fetchBansMinutes)
	if fetchBansMinutes == "" {
		serverConfig.ServerPrefs.FetchBansMinutes = defaultFetchBansMinutes
	} else {
		serverConfig.ServerPrefs.FetchBansMinutes, _ = strconv.Atoi(fetchBansMinutes)
	}

	writeConfigFile()
	fmt.Println("Config file written to : " + configPath + ", please add paths to your banlist file (and HTTPS cert/key if needed) and check over the settings.")

	fmt.Println("Press enter to exit.")
	var exit string
	fmt.Scanln(&exit)

	os.Exit(1)

}
