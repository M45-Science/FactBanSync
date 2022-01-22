package main

//Our config data
type serverConfigData struct {
	Version string
	Name    string
	ListURL string

	PathData    filePathData
	ServerPrefs serverPrefs
	WebServer   webServerConfigData
}

type serverPrefs struct {
	AutoSubscribe bool
	RequireReason bool
	StripReasons  bool

	//RCONEnabled         bool
	//LogMonitoring       bool
	//RequireMultipleBans bool

	MaxBanOutputSize int

	FetchBansMinutes int
	WatchFileSeconds int
	RefreshListHours int

	DownloadTimeoutSeconds int
	DownloadSizeLimitBytes int64
}

type filePathData struct {
	FactorioBanFile   string
	FactorioWhitelist string
	ServerListFile    string
	CompositeBanFile  string
	LogDir            string
	BanCacheDir       string
}

type webServerConfigData struct {
	RunWebServer         bool
	DomainName           string
	SSLWebPort           int
	SSLKeyFile           string
	SSLCertFile          string
	MaxRequestsPerSecond int
}

//List of servers
type serverListData struct {
	Version    string
	ServerList []serverData
}

//Server data
type serverData struct {
	Name         string
	Bans         string
	Trusts       string `json:",omitempty"`
	Logs         string `json:",omitempty"`
	Website      string `json:",omitempty"`
	Discord      string `json:",omitempty"`
	JsonGzip     bool
	UseRedScrape bool

	LocalData localData
}

type localData struct {
	Subscribed   bool
	StripReasons bool
	Added        string
	BanList      []banDataType `json:"-"`
}

//Minimal ban data
type minBanDataType struct {
	UserName string `json:"username"`
	Reason   string `json:"reason"`
}

//Ban data
type banDataType struct {
	UserName string `json:"username"`
	Reason   string `json:"reason,omitempty"`
	Revoked  bool   `json:",omitempty"`
	Added    string `json:",omitempty"`

	Sources []string `json:",omitempty"`
	Reasons []string `json:",omitempty"`
	Revokes []bool   `json:",omitempty"`
	Adds    []string `json:",omitempty"`
}

//RCON list
type RCONDataList struct {
	RCONData []RCONData
}

//RCON data
type RCONData struct {
	RCONName     string
	RCONAddress  string
	RCONPassword string
}

//Log monitor data
type LogMonitorData struct {
	Name string
	File string
	Path string
}
