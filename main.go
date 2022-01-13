package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "0.0.1"

var defaultListURL = "https://raw.githubusercontent.com/Distortions81/FactBanSync/master/server-list.json"
var defaultConfigPath = "server-config.json"
var defaultBanFile = "server-banlist.json"
var defaultServerListFile = "server-list.json"
var defaultLogPath = "logs"
var defualtFetchRate = 300
var defualtWatchInterval = 5

type serverConfigData struct {
	Version        string
	ServerName     string
	ServerURL      string
	ListURL        string
	BanFile        string
	ServerListFile string
	LogPath        string
	FetchRate      int
	WatchInterval  int
}

type serverListData struct {
	Version    string
	ServerList []serverData
}

type serverData struct {
	Subscribed bool
	ServerName string
	ServerURL  string
}

type banDataData struct {
	UserName string `json:"username"`
	Reason   string `json:"reason,omitempty"`
	Address  string `json:"address,omitempty"`
}

var serverConfig serverConfigData
var serverList serverListData
var configPath string
var logDesc *os.File
var banData []banDataData

func main() {

	configPath = *flag.String("configPath", defaultConfigPath, "config file path")
	flag.Parse()

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
	writeConfigFile()

	//Make log dir
	err = os.Mkdir(serverConfig.LogPath, 0777)

	if os.IsNotExist(err) {
		panic(err)
	}

	//Open log file
	logName := time.Now().Format("2006-01-02") + ".log"
	logDesc, err = os.OpenFile(serverConfig.LogPath+"/"+logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	defer logDesc.Close()
	mw := io.MultiWriter(os.Stdout, logDesc) //To log and stdout
	log.SetOutput(mw)

	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.Println(fmt.Sprintf("FactBanSync v%v", version))

	//Update server list
	log.Println("Updating server list")
	var wrote int64
	wrote, err = downloadFile(serverConfig.ServerListFile+".tmp", defaultListURL)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	if wrote > 0 {
		var sList serverListData
		data, err := ioutil.ReadFile(serverConfig.ServerListFile + ".tmp")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(data), &sList)
		if err == nil {
			err = os.Rename(serverConfig.ServerListFile+".tmp", serverConfig.ServerListFile)
			if err != nil {
				panic(err)
			}
			serverList = sList
		} else {
			log.Println("Unable to parse remote server list file")
			os.Remove(serverConfig.ServerListFile + ".tmp")
		}
	}
	//Read server list file
	file, err = ioutil.ReadFile(serverConfig.ServerListFile)

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

	readServerBanList()
	writeBanListFile() //To clean up formatting
}

func downloadFile(filepath string, url string) (int64, error) {

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	var wrote int64
	wrote, err = io.Copy(out, resp.Body)
	return wrote, err
}
