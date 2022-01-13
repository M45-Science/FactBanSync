package main

import "os"

var defaultListURL = "https://raw.githubusercontent.com/Distortions81/FactBanSync/master/server-list.json"
var defaultConfigPath = "server-config.json"
var defaultBanFile = "server-banlist.json"
var defaultServerListFile = "server-list.json"
var defaultOurBansFile = "our-bans.json"
var defaultLogPath = "logs"
var defualtFetchBansInterval = 60        //Minutes
var defualtWatchInterval = 5             //Seconds
var defualtRefreshListInterval = 60 * 24 //Minutes
var serverConfig serverConfigData
var serverList serverListData
var configPath string
var logDesc *os.File
var banData []banDataData
