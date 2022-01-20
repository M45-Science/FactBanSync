package main

//Our config data
type serverConfigData struct {
	Version string
	ListURL string

	Name             string
	FactorioBanFile  string
	ServerListFile   string
	CompositeBanFile string
	LogDir           string
	BanCacheDir      string
	MaxBanlistSize   int

	RunWebServer      bool
	WebPort           int
	GetTimeoutSeconds int
	GetSizeLimitBytes int64

	AutoSubscribe bool
	//RCONEnabled         bool
	//LogMonitoring       bool
	RequireReason bool
	//RequireMultipleBans bool
	StripReasons bool

	FetchBansSeconds   int
	WatchFileSeconds   int
	RefreshListMinutes int
}

//List of servers
type serverListData struct {
	Version    string
	ServerList []serverData
}

//Server data
type serverData struct {
	Name     string
	Bans     string
	Trusts   string `json:",omitempty"`
	Logs     string `json:",omitempty"`
	Website  string `json:",omitempty"`
	Discord  string `json:",omitempty"`
	JsonGzip bool

	Subscribed   bool
	UseRedScrape bool
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
