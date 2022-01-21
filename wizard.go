package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func setupWizard() {
	makeDefaultConfigFile()
	makeHTTPs := false

	fmt.Println("Pressing enter on a question will accept the default value.")
	fmt.Println("Community or server name: (Used to skip ourselves in the server list)")

	var communityName string
	fmt.Scanln(&communityName)
	if communityName == "" {
		communityName = defaultName
	}
	serverConfig.Name = communityName

	fmt.Println("Run a HTTPS (SSL) webserver, to provide server-banlist.json? (Y/n)")

	var runSSLWebServer string
	fmt.Scanln(&runSSLWebServer)
	if runSSLWebServer == "" || runSSLWebServer == "y" || runSSLWebServer == "Y" {
		serverConfig.RunWebServer = true

		fmt.Println("HTTPS Web server port (8443): ")

		var webPort string
		fmt.Scanln(&webPort)
		if webPort == "" {
			serverConfig.SSLWebPort = defaultSSLWebPort
		} else {
			serverConfig.SSLWebPort, _ = strconv.Atoi(webPort)
		}

		fmt.Println("You will need a certificate and key file for HTTPS. Put them in the data directory and put the paths in the config file. On most systems you can use the provided make-https-cert.sh script to generate self-signed certificates.")
		fmt.Println("Would you like to (attempt) to auto-run the script at the end of the setup? (y/N)")

		var runMakeHttpsCert string
		fmt.Scanln(&runMakeHttpsCert)
		if runMakeHttpsCert == "Y" || runMakeHttpsCert == "y" {
			makeHTTPs = true
		}

	} else {
		serverConfig.SSLWebPort = 0
		serverConfig.RunWebServer = false
	}

	fmt.Println("Run a HTTP (standard) webserver, to provide server-banlist.json? (y/N)")

	var runWebServer string
	fmt.Scanln(&runWebServer)
	if runWebServer == "Y" || runWebServer == "y" {
		serverConfig.RunWebServer = true

		fmt.Println("HTTP Web server port (8008): ")

		var webPort string
		fmt.Scanln(&webPort)
		if webPort == "" {
			serverConfig.WebPort = defaultWebPort
		} else {
			serverConfig.WebPort, _ = strconv.Atoi(webPort)
		}
	} else {
		serverConfig.WebPort = 0
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

	fmt.Println("How often do you want to refresh the list of servers (in hours)? (12)")

	var refreshListHours string
	fmt.Scanln(&refreshListHours)
	if refreshListHours == "" {
		serverConfig.RefreshListHours = defaultRefreshListHours
	} else {
		serverConfig.RefreshListHours, _ = strconv.Atoi(refreshListHours)
	}

	fmt.Println("How often do you want to fetch ban lists from server (in minutes)? (15)")

	var fetchBansMinutes string
	fmt.Scanln(&fetchBansMinutes)
	if fetchBansMinutes == "" {
		serverConfig.FetchBansMinutes = defaultFetchBansMinutes
	} else {
		serverConfig.FetchBansMinutes, _ = strconv.Atoi(fetchBansMinutes)
	}

	writeConfigFile()
	fmt.Println("Config file written to : " + configPath + ", please add paths to your banlist file and check over the settings.")

	if makeHTTPs {
		fmt.Println("Running make-https-cert.sh script...")

		path, _ := os.Getwd()
		os.Chdir(path)
		cmd := exec.Command("/bin/bash", "make-https-cert.sh")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run() // add error checking
	}
	os.Exit(1)

}
