package main

import (
	"io/fs"
	"sync"
	"time"
)

const ProgVersion string = "0.0.210"

// Globals
var (
	serverRunning    = true
	serverConfig     serverConfigData
	serverList       serverListData
	configPath       string
	ourBanData       []banDataType
	compositeBanData []minBanDataType

	cachedBanListGz []byte
	cachedBanList   []byte

	cachedWelcome []byte

	cachedCompositeGz   []byte
	cachedCompositeList []byte
	cachedBanListLock   sync.Mutex
	initialStat         fs.FileInfo
)

const (
	timeFormat           = time.RFC822Z
	defaultListURL       = "https://raw.githubusercontent.com/M45-Science/Factorio-Community-List/main/server-list.json"
	defaultCommunityName = "Default"

	defaultSSLWebPort             = 8443
	defaultDownloadSizeLimitKB    = 1024 //1MB
	defaultDownloadTimeoutSeconds = 30
	defaultMaxReqestsPerSecond    = 10
	defaultRunWebServer           = false
	defaultAutoSubscribe          = true
	defaultRequireReason          = false

	//PathData
	defaultBanFile = "../factorio/server-banlist.changeme"

	defaultConfigPath     = "data/server-config.json"
	defaultServerListFile = "data/server-list.json"
	defaultCompositeFile  = "data/composite.json"
	defaultFileWebName    = "server-banlist.json"
	defaultCompositeName  = "composite.json"

	defaultSSLKeyFile  = "data/server.key"
	defaultSSLCertFile = "data/server.crt"

	defaultDataDir     = "data"
	defaultLogDir      = "data/logs"
	defaultBanFileDir  = "data/banCache"
	defaultWelcomeFile = "data/welcome.html"

	//Default delay times
	defaultFetchBansMinutes = 15
	defaultWatchSeconds     = 5
	defaultRefreshListHours = 12

	//Max banlist size
	defaultMaxBanOutputCount = 10000

	defaultVerboseLogging = false
)
