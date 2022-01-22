package main

import (
	"io/fs"
	"sync"
	"time"
)

const Version string = "0.0.2"

//Globals
var (
	serverRunning     = true
	serverConfig      serverConfigData
	serverList        serverListData
	configPath        string
	ourBanData        []banDataType
	compositeBanData  []minBanDataType
	cachedBanListGz   []byte
	cachedBanList     []byte
	cachedBanListLock sync.Mutex
	initialStat       fs.FileInfo
)

const (
	defaultListURL         = "https://raw.githubusercontent.com/Distortions81/Factorio-Community-List/main/server-list.json"
	timeFormat             = time.RFC822Z
	defaultName            = "Default"
	defaultBanFile         = ""
	defaultSSLWebPort      = 8443
	defaultMaxDownloadSize = 1024 * 1024 //mb
	defaultDownloadTimeout = 10

	//Default file paths
	defaultConfigPath     = "data/server-config.json"
	defaultServerListFile = "data/server-list.json"
	defaultCompositeFile  = "data/composite.json"
	defaultFileWebName    = "server-banlist.json"
	defaultSSLKeyFile     = "data/server.key"
	defaultSSLCertFile    = "data/server.crt"

	//Default directories
	defaultDataDir    = "data"
	defaultLogDir     = "data/logs"
	defaultBanFileDir = "data/banCache"

	//Default delay times
	defaultFetchBansMinutes = 15
	defaultWatchSeconds     = 10
	defaultRefreshListHours = 12

	//Max banlist size
	defaultMaxBanListSize = 950
)
