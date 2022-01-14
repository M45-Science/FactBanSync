package main

import "time"

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

	FetchBansInterval   int
	WatchInterval       int
	RefreshListInterval int
}

//List of servers
type serverListData struct {
	Version    string
	ServerList []serverData
}

//Server data
type serverData struct {
	Subscribed   bool
	ServerName   string
	ServerURL    string
	JsonGz       bool `json:"omitempty"`
	AddedLocally time.Time
}

//Ban data
type banDataData struct {
	UserName string `json:"username"`
	Reason   string `json:"reason,omitempty"`
	Address  string `json:"address,omitempty"`
	Added    time.Time
}

//RCON list
type RCONDataList struct {
	RCONData []RCONData
}

//RCON data
type RCONData struct {
	RCONAddress  string
	RCONPassword string
}

//Log monitor data
type LogMonitorData struct {
	Name string
	Path string
}
