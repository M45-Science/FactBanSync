package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

//Read list of servers from file
func readServerListFile() {
	file, err := ioutil.ReadFile(serverConfig.ServerListFile)

	//Read server list file if it exists
	if file != nil && !os.IsNotExist(err) {
		err = json.Unmarshal(file, &serverList)

		if err != nil {
			log.Println("Error reading server list file: " + err.Error())
			os.Exit(1)
		}
	} else {
		//Generate empty list
		serverList = serverListData{Version: "0.0.1", ServerList: []serverData{}}

		log.Println("No server list file found, creating new one.")
		writeServerListFile()
	}
}

//Read server config from file
func readConfigFile() {
	//Read server config file
	file, err := ioutil.ReadFile(configPath)

	if file != nil && err == nil {
		err = json.Unmarshal(file, &serverConfig)

		if err != nil {
			log.Println("Error reading config file: " + err.Error())
			os.Exit(1)
		}

		//Let user know further config is required
		if serverConfig.Name == "Default" {
			log.Println("Please change ServerName in the config file, or use --runWizard")
			os.Exit(1)
		}
	} else {
		//Make example config file, with reasonable defaults
		fmt.Println("No config file found, generating defaults, saving to " + configPath)
		os.Mkdir(defaultDataDir, 0755)
		makeDefaultConfigFile()

		fmt.Println("Would you like to use the setup wizard? (y/N)")

		var input string
		fmt.Scanln(&input)

		if input == "y" || input == "Y" {
			setupWizard()
			return
		}

		log.Println("Exiting...")
		os.Exit(1)
	}
}

//Make default-value config file as an example starting point
func makeDefaultConfigFile() {
	serverConfig.Version = version
	serverConfig.ListURL = defaultListURL

	serverConfig.Name = defaultName
	serverConfig.FactorioBanFile = defaultBanFile
	serverConfig.ServerListFile = defaultServerListFile
	serverConfig.LogDir = defaultLogDir
	serverConfig.BanCacheDir = defaultBanFileDir
	serverConfig.MaxBanlistSize = defaultMaxBanListSize

	serverConfig.RunWebServer = false
	serverConfig.WebPort = defaultWebPort

	serverConfig.RCONEnabled = false
	serverConfig.LogMonitoring = false
	serverConfig.AutoSubscribe = false
	serverConfig.RequireReason = false

	serverConfig.FetchBansSeconds = defaultFetchBansSeconds
	serverConfig.WatchFileSeconds = defaultWatchSeconds
	serverConfig.RefreshListMinutes = defaultRefreshListMinutes

	writeConfigFile()
}

//Read the Factorio ban list file locally
func readServerBanList() {

	file, err := os.Open(serverConfig.FactorioBanFile)

	if err != nil {
		log.Println(err)
		return
	}

	var bData []banDataType

	data, err := ioutil.ReadAll(file)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var names []string
	err = json.Unmarshal(data, &names)

	if err != nil {
		//Not really an error, just empty array
		//Only needed because Factorio will write some bans as an array for some unknown reason.
	} else {

		for _, name := range names {
			if name != "" {
				bData = append(bData, banDataType{UserName: name})
			}
		}
	}

	var bans []banDataType
	err = json.Unmarshal(data, &bans)

	if err != nil {
		//Ignore, just array of strings
	}

	for _, item := range bans {
		if item.UserName != "" {
			//It also commonly writes this address, and it isn't neeeded
			if item.Address == "0.0.0.0" {
				item.Address = ""
			}
			bData = append(bData, item)
		}
	}

	ourBanData = bData

	log.Println("Read " + fmt.Sprintf("%v", len(bData)) + " bans from banlist")
	updateWebCache()
	compositeBans()
}

func updateWebCache() {

	var localCopy []banDataType
	for _, item := range ourBanData {
		if item.UserName != "" {
			var name, addr, reason string

			name = item.UserName
			if !serverConfig.StripAddresses {
				addr = item.Address
			}
			if !serverConfig.StripReasons {
				reason = item.Reason
			}

			localCopy = append(localCopy, banDataType{UserName: name, Address: addr, Reason: reason})
		}
	}

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err := enc.Encode(localCopy)

	if err != nil {
		log.Println("Error encoding ban list file: " + err.Error())
		os.Exit(1)
	}

	//Cache a normal and gzip version of the ban list
	if serverConfig.RunWebServer {
		cachedBanListLock.Lock()
		cachedBanList = outbuf.Bytes()
		cachedBanListGz = compressGzip(outbuf.Bytes())
		log.Println("Cached response: " + fmt.Sprintf("%v", len(cachedBanList)) + " json / " + fmt.Sprintf("%v", len(cachedBanListGz)) + " gz")
		cachedBanListLock.Unlock()
	}
}

//Write our ban list to the Factorio ban list file (indent)
func writeBanListFile() {
	file, err := os.Create(serverConfig.FactorioBanFile)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(ourBanData)

	if err != nil {
		log.Println("Error encoding ban list file: " + err.Error())
		os.Exit(1)
	}

	wrote, err := file.Write(outbuf.Bytes())

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Wrote banlist of " + fmt.Sprintf("%v", len(ourBanData)) + " items, " + fmt.Sprintf("%v", wrote) + " bytes")
}

//Write our server list to the server list file (indent)
func writeConfigFile() {
	file, err := os.Create(configPath)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	serverConfig.Version = "0.0.1"
	//Add config file comments

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(serverConfig)

	if err != nil {
		log.Println("Error writing config file: " + err.Error())
		os.Exit(1)
	}

	wrote, err := file.Write(outbuf.Bytes())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Wrote config file: " + fmt.Sprintf("%v", wrote) + " bytes")

}

//Write list of servers to file
func writeServerListFile() {
	file, err := os.Create(serverConfig.ServerListFile)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	serverList.Version = "0.0.1"
	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(serverList)

	if err != nil {
		log.Println("Error writing server list file: " + err.Error())
	}

	wrote, err := file.Write(outbuf.Bytes())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Print("Wrote server list file: " + fmt.Sprintf("%v", wrote) + " bytes")
}

//Sanitize a string before use in a filename
func FileNameFilter(str string) string {
	alphafilter, _ := regexp.Compile("[^a-zA-Z0-9-_]+")
	str = alphafilter.ReplaceAllString(str, "")
	return str
}

//Gzip compress a byte array
func compressGzip(data []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(data)
	w.Close()
	return b.Bytes()
}
