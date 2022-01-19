package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

/* Lots of replicated code, but I'd rather have the duplicate code than a single overly complex function */

//Save a banlist from a specific server
func saveBanLists() {
	os.Mkdir(serverConfig.BanCacheDir, 0777)
	for _, server := range serverList.ServerList {
		if server.Name == server.Name {
			continue
		}
		if server.Subscribed {
			log.Println("Saving ban list for server: " + server.Name + " (" + strconv.Itoa(len(server.BanList)) + " bans)")

			outbuf := new(bytes.Buffer)
			enc := json.NewEncoder(outbuf)
			enc.SetIndent("", "\t")

			err := enc.Encode(server.BanList)

			if err != nil {
				log.Println("Error encoding ban list file: " + err.Error())
				os.Exit(1)
			}
			err = ioutil.WriteFile(defaultBanFileDir+"/"+FileNameFilter(server.Name)+".json", outbuf.Bytes(), 0644)
			if err != nil {
				log.Println("Error saving ban list: " + err.Error())
				continue
			}
		}
	}
}

//Make default-value config file as an example starting point
func makeDefaultConfigFile() {
	serverConfig.Version = version
	serverConfig.ListURL = defaultListURL

	serverConfig.Name = defaultName
	serverConfig.FactorioBanFile = defaultBanFile
	serverConfig.ServerListFile = defaultServerListFile
	serverConfig.LogDir = defaultLogDir
	serverConfig.BanCacheDir = defaultBanFileDir
	serverConfig.MaxBanlistSize = defaultMaxBanListSize

	serverConfig.RunWebServer = false
	serverConfig.WebPort = defaultWebPort

	//serverConfig.RCONEnabled = false
	//serverConfig.LogMonitoring = false
	serverConfig.AutoSubscribe = true
	serverConfig.RequireReason = false

	serverConfig.FetchBansSeconds = defaultFetchBansSeconds
	serverConfig.WatchFileSeconds = defaultWatchSeconds
	serverConfig.RefreshListMinutes = defaultRefreshListMinutes

	writeConfigFile()
}

//Write our ban list to the Factorio ban list file
func writeBanListFile() {
	file, err := os.Create(serverConfig.FactorioBanFile)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(ourBanData)

	if err != nil {
		log.Println("Error encoding ban list file: " + err.Error())
		os.Exit(1)
	}

	wrote, err := file.Write(outbuf.Bytes())

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Wrote banlist of " + fmt.Sprintf("%v", len(ourBanData)) + " items, " + fmt.Sprintf("%v", wrote) + " bytes")
}

//Write our server list to the server list file (indent)
func writeConfigFile() {
	file, err := os.Create(configPath)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	serverConfig.Version = "0.0.1"
	//Add config file comments

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

	log.Println("Wrote config file: " + fmt.Sprintf("%v", wrote) + "b")

}

//Write list of servers to file
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

//Write out a combined list of bans
func writeCompositeBanlist() {
	file, err := os.Create(serverConfig.CompositeBanFile)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(compositeBanData)

	if err != nil {
		log.Println("writeCompositeBanlist: " + err.Error())
		os.Exit(1)
	}

	wrote, err := file.Write(outbuf.Bytes())

	if err != nil {
		log.Println("writeCompositeBanlist: " + err.Error())
	}

	log.Println("Wrote composite banlist of " + fmt.Sprintf("%v", len(compositeBanData)) + " items, " + fmt.Sprintf("%v", wrote/1024) + " kb")
}
