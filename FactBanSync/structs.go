package main

//Our config data
type serverConfigData struct {
	Version string
	ListURL string

	ServerName     string
	BanFile        string
	ServerListFile string
	LogDir         string
	BanFileDir     string

	RunWebServer bool
	WebPort      int

	RCONEnabled   bool
	LogMonitoring bool
	AutoSubscribe bool
	RequireReason bool

	FetchBansSeconds   int
	WatchSeconds       int
	RefreshListMinutes int
}

//List of servers
type serverListData struct {
	Version    string
	ServerList []serverData
}

//Server data
type serverData struct {
	ServerName string
	ServerURL  string
	Website    string `json:"omitempty"`
	Discord    string `json:"omitempty"`
	Logs       string `json:"omitempty"`
	JsonGzip   bool   `json:"omitempty"`
	Subscribed bool   `json:"omitempty"`
	LocalAdd   string

	BanList []banDataType `json:"-"`
}

//Ban data
type banDataType struct {
	UserName string `json:"username"`
	Reason   string `json:"reason,omitempty"`
	Address  string `json:"address,omitempty"`
	LocalAdd string `json:"omitempty"`
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
