package main

import "time"

type serverConfigData struct {
	Version string

	Comment1   string
	ServerName string

	Comment2 string
	ListURL  string

	Comment3 string
	BanFile  string

	Comment4       string
	ServerListFile string

	Comment5 string
	LogPath  string

	Comment6      string
	AutoSubscribe bool

	Comment7      string
	RequireReason bool

	Comment8          string
	FetchBansInterval int

	Comment9      string
	WatchInterval int

	Comment10           string
	RefreshListInterval int

	Comment11   string
	OurBansFile string
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
