package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

/* Lots of dupe code, but I'd rather have the duplicate code than a single overly complex function */

func readBanCache() {
	for spos, server := range serverList.ServerList {
		if strings.EqualFold(server.CommunityName, serverConfig.CommunityName) {
			continue
		}
		serverList.ServerList[spos].LocalData.BanList = readBanCacheFile(server.CommunityName)
	}
}

func readBanCacheFile(serverName string) []banDataType {

	bandata := []banDataType{}
	serverName = FileNameFilter(serverName)
	path := serverConfig.PathData.BanCacheDir + "/" + serverName + ".json"
	file, err := os.ReadFile(path)

	if file != nil && err == nil {
		err = json.Unmarshal(file, &bandata)

		if err != nil {
			log.Println("Error reading ban cache file: " + err.Error())
		}
	}

	return bandata

}

// Read list of servers from file
func readServerListFile() {
	file, err := os.ReadFile(serverConfig.PathData.ServerListFile)

	//Read server list file if it exists
	if file != nil && !os.IsNotExist(err) {
		var temp serverListData
		err = json.Unmarshal(file, &temp)

		if err != nil {
			log.Println("Error reading server list file: " + err.Error())
			os.Exit(1)
		}

		if temp.Version == "0.0.1" {
			serverList = temp
		} else {
			log.Print("Server list file version is not supported, exiting.")
			time.Sleep(time.Minute)
			os.Exit(1)
		}
	} else {
		//Generate empty list
		serverList = serverListData{Version: "0.0.1", ServerList: []serverData{}}

		log.Println("No server list file found, creating new one.")
		writeServerListFile()
	}
}

// Read server config from file
func readConfigFile() {
	//Read server config file
	file, err := os.ReadFile(configPath)

	if file != nil && err == nil {
		var temp serverConfigData
		err = json.Unmarshal(file, &temp)

		if err != nil {
			log.Println("Error reading config file: " + err.Error())
			os.Exit(1)
		}

		//Let user know further config is required
		if strings.EqualFold(serverConfig.CommunityName, "Default") {
			log.Println("Please change ServerName in the config file, or use --runWizard")
			os.Exit(1)
		}

		if temp.Version == "0.0.2" {
			serverConfig = temp
		} else {
			log.Print("Config file version is not supported, exiting.")
			time.Sleep(time.Minute)
			os.Exit(1)
		}
	} else {
		//Make example config file, with reasonable defaults
		log.Println("No config file found, generating defaults and saving to " + configPath)
		os.Mkdir(defaultDataDir, 0755)
		makeDefaultConfigFile()

		fmt.Println("Would you like to use the setup wizard? (Y/n)")

		var input string
		fmt.Scanln(&input)

		if input == "y" || input == "Y" || input == "" {
			setupWizard()
			return
		} else {
			log.Println("Please edit the config file, or use --runWizard")
		}

		log.Println("Exiting...")
		os.Exit(1)
	}
}

// Read the Factorio ban list file locally
func readServerBanList() {

	if serverConfig.PathData.FactorioBanFile == "" {
		return
	}
	file, err := os.Open(serverConfig.PathData.FactorioBanFile)

	if err != nil {
		log.Println(err)
		return
	}

	var bData []banDataType

	data, err := io.ReadAll(file)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var names []string
	err = json.Unmarshal(data, &names)

	if err != nil {
		//Not really an error, just empty array
		//Only needed because Factorio will write some bans as an array for some unknown reason.
	} else {

		for npos, name := range names {
			if name != "" {
				bData = append(bData, banDataType{UserName: strings.ToLower(name), Added: time.Now().Format(timeFormat)})
				names[npos] = strings.ToLower(name)
			}
		}
	}

	var bans []banDataType
	err = json.Unmarshal(data, &bans)

	if err != nil {
		fmt.Print("") //Annoying warning remover
	}

	for ipos, item := range bans {
		if item.UserName != "" && !strings.HasPrefix(strings.ToLower(item.Reason), strings.ToLower("[FCL]")) {
			bData = append(bData, banDataType{UserName: strings.ToLower(item.UserName), Reason: item.Reason, Added: time.Now().Format(timeFormat)})
			bans[ipos].UserName = strings.ToLower(item.UserName)
		}
	}

	ourBanData = bData

	if serverConfig.ServerPrefs.VerboseLogging {
		log.Printf("Read %v bans from banlist.\n", len(bData))
	}
	compositeBans()
	updateWebCache()
}
