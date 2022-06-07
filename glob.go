package main

import (
	"io/fs"
	"sync"
	"time"
)

const ProgVersion string = "0.0.202"

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
	timeFormat           = time.RFC822Z
	defaultListURL       = "https://raw.githubusercontent.com/Distortions81/Factorio-Community-List/main/server-list.json"
	defaultCommunityName = "Default"

	defaultSSLWebPort             = 8443
	defaultDownloadSizeLimitKB    = 1024 //1MB
	defaultDownloadTimeoutSeconds = 10
	defaultMaxReqestsPerSecond    = 10
	defaultRunWebServer           = false
	defaultAutoSubscribe          = true
	defaultRequireReason          = false

	//PathData
	defaultBanFile        = "../factorio/server-banlist.changeme"
	defaultConfigPath     = "data/server-config.json"
	defaultServerListFile = "data/server-list.json"
	defaultCompositeFile  = "data/composite.json"
	defaultFileWebName    = "server-banlist.json"
	defaultSSLKeyFile     = "data/server.key"
	defaultSSLCertFile    = "data/server.crt"
	defaultDataDir        = "data"
	defaultLogDir         = "data/logs"
	defaultBanFileDir     = "data/banCache"

	//Default delay times
	defaultFetchBansMinutes = 15
	defaultWatchSeconds     = 10
	defaultRefreshListHours = 12

	//Max banlist size
	defaultMaxBanOutputCount = 10000

	defaultVerboseLogging = false
)
