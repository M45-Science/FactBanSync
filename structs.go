package main

//Our config data
type serverConfigData struct {
	Version       string
	CommunityName string
	ServerListURL string

	PathData    filePathData
	ServerPrefs serverPrefs
	WebServer   webServerConfigData
}

type serverPrefs struct {
	AutoSubscribe bool
	RequireReason bool
	StripReasons  bool

	MaxBanOutputCount int
	WatchFileSeconds  int
	FetchBansMinutes  int
	RefreshListHours  int

	DownloadTimeoutSeconds int
	DownloadSizeLimitKB    int64

	VerboseLogging bool
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
	CommunityName string `json:"Name"`
	BanListURL    string `json:"Bans"`
	WhiteListURL  string `json:"Trusts,omitempty"`
	LogFileURL    string `json:"Logs,omitempty"`
	WebsiteURL    string `json:"Website,omitempty"`
	DiscordURL    string `json:"Discord,omitempty"`
	JsonGzip      bool
	UseRedScrape  bool

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
