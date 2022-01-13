package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func readServerListFile() {
	file, err := ioutil.ReadFile(serverConfig.ServerListFile)

	if file != nil && !os.IsNotExist(err) {
		err = json.Unmarshal(file, &serverList)

		if err != nil {
			log.Println("Error reading server list file: " + err.Error())
			os.Exit(1)
		}
	} else {
		serverList = serverListData{Version: "0.0.1", ServerList: []serverData{{ServerName: serverConfig.ServerName, Subscribed: true, Added: time.Now()}}}

		log.Println("No server list file found, creating new one.")
		writeServerListFile()
		os.Exit(1)
	}
}

func readConfigFile() {
	//Read server config file
	file, err := ioutil.ReadFile(configPath)

	if file != nil && err == nil {
		err = json.Unmarshal(file, &serverConfig)

		if err != nil {
			log.Println("Error reading config file: " + err.Error())
			os.Exit(1)
		}

		if serverConfig.ServerName == "Default" {
			log.Println("Please change ServerName in the config file")
			os.Exit(1)
		}
	} else {
		serverConfig.Version = version
		serverConfig.ServerName = "Default"
		serverConfig.ListURL = defaultListURL
		serverConfig.BanFile = defaultBanFile
		serverConfig.ServerListFile = defaultServerListFile
		serverConfig.LogPath = defaultLogPath
		serverConfig.AutoSubscribe = true
		serverConfig.FetchBansInterval = defualtFetchBansInterval
		serverConfig.WatchInterval = defualtWatchInterval
		serverConfig.RefreshListInterval = defualtRefreshListInterval
		serverConfig.OurBansFile = defaultOurBansFile

		fmt.Println("No config file found, generating defaults, saving to " + configPath)
		log.Println("Please change ServerName in the config file!")
		log.Println("Exiting...")
		writeConfigFile()
		os.Exit(1)
	}
}

func readServerBanList() {

	file, err := os.Open(serverConfig.BanFile)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var bData []banDataData

	data, err := ioutil.ReadAll(file)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var names []string
	err = json.Unmarshal([]byte(data), &names)

	if err != nil {
		//Not really an error, just empty array
		//Only needed because Factorio will write some bans as an array for some unknown reason.
	}

	for _, name := range names {
		if name != "" {
			bData = append(bData, banDataData{UserName: name})
		}
	}

	var bans []banDataData
	err = json.Unmarshal([]byte(data), &bans)

	if err != nil {
		log.Println("Error reading ban list file: " + err.Error())
		os.Exit(1)
	}

	for _, item := range bans {
		if item.UserName != "" {
			if item.Address == "0.0.0.0" {
				item.Address = ""
			}
			bData = append(bData, item)
		}
	}

	banData = bData

	log.Println("Read " + fmt.Sprintf("%v", len(bData)) + " bans from banlist")
}

func writeBanListFile() {
	file, err := os.Create(serverConfig.BanFile)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(banData)

	if err != nil {
		log.Println("Error writing ban list file: " + err.Error())
		os.Exit(1)
	}

	wrote, err := file.Write(outbuf.Bytes())

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Wrote banlist of " + fmt.Sprintf("%v", len(banData)) + " items, " + fmt.Sprintf("%v", wrote) + " bytes")
}

func writeConfigFile() {
	file, err := os.Create(configPath)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	serverConfig.Version = "0.0.1"
	serverConfig.Comment1 = "Your server name. This is used to skip your own server in the server list."
	serverConfig.Comment2 = "URL where server list is located, normally the git repo."
	serverConfig.Comment3 = "Path to your Factorio server-banlist.json."
	serverConfig.Comment4 = "Path to store server list locally. Cache, and allows manual subscription to servers."
	serverConfig.Comment5 = "Path of directory to put our logs."
	serverConfig.Comment6 = "Auto-subscribe to new servers."
	serverConfig.Comment7 = "Only accept bans with a reason specified."
	serverConfig.Comment8 = "How often to check other servers for new bans (in minutes)."
	serverConfig.Comment9 = "How often to check for new bans on our own server. (seconds)"
	serverConfig.Comment10 = "How often to check for new servers. (minutes)"

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(serverConfig)

	if err != nil {
		log.Println("Error writing config file: " + err.Error())
		os.Exit(1)
	}

	wrote, err := file.Write(outbuf.Bytes())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Wrote config file: " + fmt.Sprintf("%v", wrote) + " bytes")

}

func writeServerListFile() {
	file, err := os.Create(serverConfig.ServerListFile)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	serverList.Version = "0.0.1"
	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(serverList)

	if err != nil {
		log.Println("Error writing server list file: " + err.Error())
	}

	wrote, err := file.Write(outbuf.Bytes())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Print("Wrote server list file: " + fmt.Sprintf("%v", wrote) + " bytes")
}

func readServerList() {

	file, err := os.Open(serverConfig.ServerListFile)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var sList serverListData

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(data), &sList)

	if err != nil {
		log.Println("Error reading server list file: " + err.Error())
		os.Exit(1)
	}

	serverList = sList
}
