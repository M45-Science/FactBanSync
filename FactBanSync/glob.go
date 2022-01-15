package main

import (
	"io/fs"
	"os"
	"sync"
	"time"
)

var serverRunning = true
var defaultListURL = "https://raw.githubusercontent.com/Distortions81/Factorio-Community-List/main/server-list.json"
var timeFormat = time.RFC822Z

//Default file names
var defaultBanFile = ""
var defaultConfigPath = "data/server-config.json"
var defaultServerListFile = "data/server-list.json"
var defaultRCONFile = "data/server-rcon.json"
var defaultLogMonitorFile = "data/log-monitor.json"
var banFileWebName = "server-banlist.json"

//Default directories
var defaultLogDir = "data/logs"
var defaultBanFileDir = "data/banCache"

//Default delay times
var defaultFetchBansSeconds = 15        //Seconds
var defaultWatchSeconds = 10            //Seconds
var defaultRefreshListMinutes = 60 * 24 //One a day

//Global vars
var serverConfig serverConfigData
var serverList serverListData
var configPath string
var logDesc *os.File
var banData []banDataType
var logsToMonitor []LogMonitorData
var cachedBanListGz []byte
var cachedBanList []byte
var cachedBanListLock sync.Mutex

var initialStat fs.FileInfo
