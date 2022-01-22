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
	os.Mkdir(serverConfig.PathData.BanCacheDir, 0777)
	for _, server := range serverList.ServerList {
		if server.CommunityName == serverConfig.CommunityName {
			continue
		}
		if server.LocalData.Subscribed {
			if serverConfig.ServerPrefs.VerboseLogging {
				log.Println("Saving ban list for server: " + server.CommunityName + " (" + strconv.Itoa(len(server.LocalData.BanList)) + " bans)")
			}

			outbuf := new(bytes.Buffer)
			enc := json.NewEncoder(outbuf)
			enc.SetIndent("", "\t")

			err := enc.Encode(server.LocalData.BanList)

			if err != nil {
				log.Println("Error encoding ban list file: " + err.Error())
				os.Exit(1)
			}
			err = ioutil.WriteFile(defaultBanFileDir+"/"+FileNameFilter(server.CommunityName)+".json", outbuf.Bytes(), 0644)
			if err != nil {
				if serverConfig.ServerPrefs.VerboseLogging {
					log.Println("Error saving ban list: " + err.Error())
				}
				continue
			}
		}
	}
}

//Make default-value config file as an example starting point
func makeDefaultConfigFile() {
	serverConfig.ServerListURL = defaultListURL

	serverConfig.CommunityName = defaultCommunityName

	serverConfig.PathData.FactorioBanFile = defaultBanFile
	serverConfig.PathData.ServerListFile = defaultServerListFile
	serverConfig.PathData.LogDir = defaultLogDir
	serverConfig.PathData.BanCacheDir = defaultBanFileDir
	serverConfig.PathData.CompositeBanFile = defaultCompositeFile
	serverConfig.PathData.FactorioBanFile = defaultBanFile

	serverConfig.WebServer.RunWebServer = defaultRunWebServer
	serverConfig.WebServer.SSLKeyFile = defaultSSLKeyFile
	serverConfig.WebServer.SSLCertFile = defaultSSLCertFile
	serverConfig.WebServer.SSLWebPort = defaultSSLWebPort
	serverConfig.WebServer.MaxRequestsPerSecond = defaultMaxReqestsPerSecond

	serverConfig.ServerPrefs.AutoSubscribe = defaultAutoSubscribe
	serverConfig.ServerPrefs.RequireReason = defaultRequireReason
	serverConfig.ServerPrefs.MaxBanOutputSize = defaultMaxBanListSize
	serverConfig.ServerPrefs.FetchBansMinutes = defaultFetchBansMinutes
	serverConfig.ServerPrefs.WatchFileSeconds = defaultWatchSeconds
	serverConfig.ServerPrefs.RefreshListHours = defaultRefreshListHours
	serverConfig.ServerPrefs.DownloadSizeLimitKB = defaultDownloadSizeLimitKB
	serverConfig.ServerPrefs.DownloadTimeoutSeconds = defaultDownloadTimeoutSeconds
	serverConfig.ServerPrefs.VerboseLogging = defaultVerboseLogging

	writeConfigFile()
}

//Write our server list to the server list file (indent)
func writeConfigFile() {
	file, err := os.Create(configPath)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	serverConfig.Version = "0.0.2"
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

	if serverConfig.ServerPrefs.VerboseLogging {
		log.Println("Wrote config file: " + fmt.Sprintf("%v", wrote) + "b")
	}

}

//Write list of servers to file
func writeServerListFile() {
	file, err := os.Create(serverConfig.PathData.ServerListFile)

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
	if serverConfig.ServerPrefs.VerboseLogging {
		log.Print("Wrote server list file: " + fmt.Sprintf("%v", wrote) + " bytes")
	}
}

//Write out a combined list of bans
func writeCompositeBanlist() {

	file, err := os.Create(serverConfig.PathData.CompositeBanFile)

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

	if serverConfig.ServerPrefs.VerboseLogging {
		log.Println("Wrote composite banlist of " + fmt.Sprintf("%v", len(compositeBanData)) + " items, " + fmt.Sprintf("%v", wrote/1024) + " kb")
	}
}
