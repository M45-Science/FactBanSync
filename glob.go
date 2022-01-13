package main

import "os"

var defaultListURL = "https://raw.githubusercontent.com/Distortions81/FactBanSync/master/server-list.json"
var defaultConfigPath = "server-config.json"
var defaultBanFile = "server-banlist.json"
var defaultServerListFile = "server-list.json"
var defaultLogPath = "logs"
var defualtFetchRate = 300
var defualtWatchInterval = 5

var serverConfig serverConfigData
var serverList serverListData
var configPath string
var logDesc *os.File
var banData []banDataData
