package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func readServerListFile() {
	file, err := ioutil.ReadFile(serverConfig.ServerListFile)

	if file != nil && !os.IsNotExist(err) {
		err = json.Unmarshal(file, &serverList)

		if err != nil {
			log.Println(err)
			panic(err)
		}
	} else {
		serverList = serverListData{Version: "0.0.1", ServerList: []serverData{{ServerName: serverConfig.ServerName, ServerURL: serverConfig.ServerURL, Subscribed: true}}}

		log.Println("No server list file found, creating new one.")
		WriteServerListFile()
		os.Exit(1)
	}
}

func readConfigFile() {
	//Read server config file
	file, err := ioutil.ReadFile(configPath)

	if file != nil && err == nil {
		err = json.Unmarshal(file, &serverConfig)

		if err != nil {
			panic(err)
		}

		if serverConfig.ServerName == "Default" {
			log.Println("Please change ServerName in the config file")
			os.Exit(1)
		}
	} else {
		serverConfig.Version = version
		serverConfig.ServerName = "Default"
		serverConfig.BanFile = defaultBanFile
		serverConfig.ServerListFile = defaultServerListFile
		serverConfig.ListURL = defaultListURL
		serverConfig.LogPath = defaultLogPath
		serverConfig.FetchRate = defualtFetchRate
		serverConfig.WatchInterval = defualtWatchInterval

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
		panic(err)
	}

	var bData []banDataData

	data, err := ioutil.ReadAll(file)

	var names []string
	_ = json.Unmarshal([]byte(data), &names)

	for _, name := range names {
		if name != "" {
			bData = append(bData, banDataData{UserName: name})
		}
	}

	var bans []banDataData
	_ = json.Unmarshal([]byte(data), &bans)

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
		panic(err)
	}

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(banData)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	_, err = file.Write(outbuf.Bytes())

	log.Println("Wrote banlist of " + fmt.Sprintf("%v", len(banData)) + " items")
}

func writeConfigFile() {
	file, err := os.Create(configPath)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(serverConfig)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	_, err = file.Write(outbuf.Bytes())

}

func WriteServerListFile() {
	file, err := os.Create(serverConfig.ServerListFile)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	serverList.Version = "0.0.1"
	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(serverList)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	_, err = file.Write(outbuf.Bytes())
	log.Print("Wrote server list file")
}

func readServerList() {

	file, err := os.Open(serverConfig.ServerListFile)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	var sList serverListData

	data, err := ioutil.ReadAll(file)

	err = json.Unmarshal([]byte(data), &sList)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	serverList = sList
}
