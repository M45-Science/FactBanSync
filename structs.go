package main

type serverConfigData struct {
	Version        string
	ServerName     string
	ServerURL      string
	ListURL        string
	BanFile        string
	ServerListFile string
	LogPath        string
	FetchRate      int
	WatchInterval  int
}

type serverListData struct {
	Version    string
	ServerList []serverData
}

type serverData struct {
	Subscribed bool
	ServerName string
	ServerURL  string
}

type banDataData struct {
	UserName string `json:"username"`
	Reason   string `json:"reason,omitempty"`
	Address  string `json:"address,omitempty"`
}
