package main

import "time"

type serverConfigData struct {
	Version string
	ListURL string

	ServerName     string
	BanFile        string
	ServerListFile string
	LogPath        string

	RCONAddresss string
	RCONPassword string
	RunWebServer bool
	WebPort      int

	RCONEnabled   bool
	LogMonitoring bool
	AutoSubscribe bool
	RequireReason bool

	FetchBansInterval   int
	WatchInterval       int
	RefreshListInterval int
}

type serverListData struct {
	Version    string
	ServerList []serverData
}

type serverData struct {
	Subscribed bool
	ServerName string
	ServerURL  string
	Added      time.Time
}

type banDataData struct {
	UserName string `json:"username"`
	Reason   string `json:"reason,omitempty"`
	Address  string `json:"address,omitempty"`
}

type RCONDataList struct {
	RCONData []RCONData
}

type RCONData struct {
	RCONAddress  string
	RCONPassword string
}

type LogMonitorData struct {
	Name string
	Path string
}
