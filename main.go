package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const version = "0.0.1"

var defaultConfigPath = "serverConfig.json"
var defaultBanFile = "server-banlist.json"
var defaultServerListFile = "serverList.json"
var defaultLogPath = "logs"
var defualtFetchRate = 300
var defualtWatchInterval = 5

type serverConfigData struct {
	Version        string
	ServerName     string
	BanFile        string
	ServerListFile string
	LogPath        string
	FetchRate      int
	URL            string
	WatchInterval  int
}

type serverListData struct {
	Version    string
	ServerList []serverData
}

type serverData struct {
	Subscribed    bool
	ServerName    string
	ServerAddress string
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
		serverConfig.LogPath = defaultLogPath
		serverConfig.FetchRate = defualtFetchRate
		serverConfig.WatchInterval = defualtWatchInterval

		fmt.Println("No config file found, generating defaults, saving to " + configPath)
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

	//Read server list file
	file, err = ioutil.ReadFile(serverConfig.ServerListFile)

	if file != nil && !os.IsNotExist(err) {
		var serverList serverListData
		err = json.Unmarshal(file, &serverList)

		if err != nil {
			log.Println(err)
			panic(err)
		}
	} else {
		var serverList serverListData
		serverList.ServerList = make([]serverData, 0)
		serverList.Version = "0.0.1"
		log.Println("No server list file found, creating new one.")
		log.Println("Please change ServerName in the config file!")
		log.Println("Exiting...")

		WriteServerListFile()
	}

	readServerBanList()
	writeBanListFile() //To clean up formatting
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
