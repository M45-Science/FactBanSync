package main

import "time"

type serverConfigData struct {
	Version             string
	ServerName          string
	ListURL             string
	BanFile             string
	ServerListFile      string
	LogPath             string
	AutoSubscribe       bool
	RequireReason       bool
	FetchBansInterval   int
	WatchInterval       int
	RefreshListInterval int
	RCONEnabled         bool
	RCONAddresss        string
	RCONPassword        string
	RunWebServer        bool
	WebPort             int
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
