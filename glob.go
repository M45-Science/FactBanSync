package main

import (
	"io/fs"
	"os"
	"sync"
	"time"
)

//Server defaults
var serverRunning = true
var defaultListURL = "https://raw.githubusercontent.com/Distortions81/Factorio-Community-List/main/server-list.json"
var timeFormat = time.RFC822Z

const defaultName = "Default"
const defaultBanFile = ""
const defaultWebPort = 8080

//Default file paths
const defaultConfigPath = "data/server-config.json"
const defaultServerListFile = "data/server-list.json"
const defaultRCONFile = "data/server-rcon.json"
const defaultLogMonitorFile = "data/log-monitor.json"
const defaultCompositeFile = "data/composite.json"
const defaultFileWebName = "server-banlist.json"

//Default directories
const defaultDataDir = "data"
const defaultLogDir = "data/logs"
const defaultBanFileDir = "data/banCache"

//Default delay times
const defaultFetchBansSeconds = 300
const defaultWatchSeconds = 10
const defaultRefreshListMinutes = 60 * 12

//Max banlist size
const defaultMaxBanListSize = 950

//Global vars
var serverConfig serverConfigData
var serverList serverListData
var configPath string
var logDesc *os.File
var ourBanData []banDataType
var compositeBanData []minBanDataType
var logsToMonitor []LogMonitorData
var cachedBanListGz []byte
var cachedBanList []byte
var cachedBanListLock sync.Mutex

var initialStat fs.FileInfo
