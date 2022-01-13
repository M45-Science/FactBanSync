package main

import (
	"io/fs"
	"os"
)

var defaultListURL = "https://raw.githubusercontent.com/Distortions81/FactBanSync/master/server-list.json"
var defaultConfigPath = "server-config.json"
var defaultBanFile = "server-banlist.json"
var defaultServerListFile = "server-list.json"
var defaultRCONFile = "server-rcon.json"
var defaultOurBansFile = "our-bans.json"
var defaultLogPath = "logs"
var defualtFetchBansInterval = 5    //Seconds
var defualtWatchInterval = 5        //Seconds
var defualtRefreshListInterval = 60 //Minutes
var serverConfig serverConfigData
var serverList serverListData
var configPath string
var logDesc *os.File
var banData []banDataData
var serverRunning = true

var initialStat fs.FileInfo
