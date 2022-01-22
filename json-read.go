package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

/* Lots of dupe code, but I'd rather have the duplicate code than a single overly complex function */

func readBanCache() {
	for spos, server := range serverList.ServerList {
		if server.CommunityName == serverConfig.CommunityName {
			continue
		}
		serverList.ServerList[spos].LocalData.BanList = readBanCacheFile(spos, server.CommunityName)
	}
}

func readBanCacheFile(spos int, serverName string) []banDataType {

	bandata := []banDataType{}
	serverName = FileNameFilter(serverName)
	path := serverConfig.PathData.BanCacheDir + "/" + serverName + ".json"
	file, err := ioutil.ReadFile(path)

	if file != nil && err == nil {
		err = json.Unmarshal(file, &bandata)

		if err != nil {
			log.Println("Error reading ban cache file: " + err.Error())
		}
	}

	return bandata

}

//Read list of servers from file
func readServerListFile() {
	file, err := ioutil.ReadFile(serverConfig.PathData.ServerListFile)

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

//Read server config from file
func readConfigFile() {
	//Read server config file
	file, err := ioutil.ReadFile(configPath)

	if file != nil && err == nil {
		var temp serverConfigData
		err = json.Unmarshal(file, &temp)

		if err != nil {
			log.Println("Error reading config file: " + err.Error())
			os.Exit(1)
		}

		//Let user know further config is required
		if serverConfig.CommunityName == "Default" {
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
		fmt.Println("No config file found, generating defaults, saving to " + configPath)
		os.Mkdir(defaultDataDir, 0755)
		makeDefaultConfigFile()

		fmt.Println("Would you like to use the setup wizard? (Y/n)")

		var input string
		fmt.Scanln(&input)

		if input == "y" || input == "Y" || input == "" {
			setupWizard()
			return
		}

		log.Println("Exiting...")
		os.Exit(1)
	}
}

//Read the Factorio ban list file locally
func readServerBanList() {

	file, err := os.Open(serverConfig.PathData.FactorioBanFile)

	if err != nil {
		log.Println(err)
		return
	}

	var bData []banDataType

	data, err := ioutil.ReadAll(file)

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

		for _, name := range names {
			if name != "" {
				bData = append(bData, banDataType{UserName: name, Added: time.Now().Format(timeFormat)})
			}
		}
	}

	var bans []banDataType
	err = json.Unmarshal(data, &bans)

	if err != nil {
		fmt.Print("") //Annoying warning remover
	}

	for _, item := range bans {
		if item.UserName != "" && !strings.HasPrefix(item.Reason, "[auto]") {
			bData = append(bData, banDataType{UserName: item.UserName, Reason: item.Reason, Added: time.Now().Format(timeFormat)})
		}
	}

	ourBanData = bData

	if serverConfig.ServerPrefs.VerboseLogging {
		log.Println("Read " + fmt.Sprintf("%v", len(bData)) + " bans from banlist")
	}
	updateWebCache()
	compositeBans()
}
