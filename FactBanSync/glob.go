package main

import (
	"io/fs"
	"os"
)

var defaultListURL = "https://raw.githubusercontent.com/Distortions81/FactBanSync/master/server-list.json"
var defaultConfigPath = "server-config.json"
var defaultBanFile = "server-banlist.json"
var defaultServerListFile = "server-list.json"
var defaultRCONFile = "server-rcon.json"
var defaultOurBansFile = "our-bans.json"
var defaultLogMonitorFile = "log-monitor.json"
var defaultLogPath = "logs"
var banFileWebName = "server-banlist.json"
var defualtFetchBansInterval = 10        //Seconds
var defualtWatchInterval = 5             //Seconds
var defualtRefreshListInterval = 60 * 24 //One a day
var serverConfig serverConfigData
var serverList serverListData
var configPath string
var logDesc *os.File
var banData []banDataData
var logsToMonitor []LogMonitorData
var cachedBanListGz []byte
var cachedBanList []byte
var serverRunning = true

var initialStat fs.FileInfo
