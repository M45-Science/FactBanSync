package main

import (
	"io/fs"
	"os"
)

var serverRunning = true
var defaultListURL = "https://raw.githubusercontent.com/Distortions81/FactBanSync/master/server-list.json"

//Default file names
var defaultConfigPath = "server-config.json"
var defaultBanFile = "server-banlist.json"
var defaultServerListFile = "server-list.json"
var defaultRCONFile = "server-rcon.json"
var defaultOurBansFile = "our-bans.json"
var defaultLogMonitorFile = "log-monitor.json"
var banFileWebName = "server-banlist.json"

//Default directories
var defaultLogDir = "logs"
var defaultBanFileDir = "banLists"

//Defualt delay times
var defualtFetchBansInterval = 15        //Seconds
var defualtWatchInterval = 10            //Seconds
var defualtRefreshListInterval = 60 * 24 //One a day

//Glboal vars
var serverConfig serverConfigData
var serverList serverListData
var configPath string
var logDesc *os.File
var banData []banDataData
var logsToMonitor []LogMonitorData
var cachedBanListGz []byte
var cachedBanList []byte

var initialStat fs.FileInfo
