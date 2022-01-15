package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func fetchBanLists() {
	gDirty := 0
	for spos, server := range serverList.ServerList {
		lDirty := 0
		if server.Subscribed {
			data, err := fetchFile(server.ServerURL)
			if err != nil {
				log.Println("Error updating ban list: " + err.Error())
				continue
			}
			if len(data) > 0 {

				var names []string
				err = json.Unmarshal(data, &names)

				if err != nil {
					//Not really an error, just empty array
					//Only needed because Factorio will write some bans as an array for some unknown reason.
				} else {

					found := false
					if len(names) > 0 {
						for _, name := range names {
							if name != "" {
								for _, item := range server.BanList {
									if item.UserName == name {
										found = true
									}
								}
								if !found {
									gDirty++
									lDirty++
									serverList.ServerList[spos].BanList = append(serverList.ServerList[spos].BanList, banDataType{UserName: name, LocalAdd: time.Now().Format(timeFormat)})
								}
							}
						}
					}
				}

				var bans []banDataType
				err = json.Unmarshal(data, &bans)

				if err != nil {
					//Ignore, just array of strings
				}

				if len(bans) > 0 {
					for _, item := range bans {
						if item.UserName != "" {
							//It also commonly writes this address, and it isn't neeeded
							if item.Address == "0.0.0.0" {
								item.Address = ""
							}
							found := false
							for _, ban := range server.BanList {
								if ban.UserName == item.UserName {
									found = true
								}
							}
							if !found {
								gDirty++
								lDirty++
								serverList.ServerList[spos].BanList = append(serverList.ServerList[spos].BanList, item)
							}
						}
					}
				}

			}
			if lDirty > 0 {
				log.Printf("Found %v new bans for %v\n", lDirty, server.ServerName)
			}
		}
	}
	if gDirty > 0 {
		saveBanLists()
	}
}

func saveBanLists() {
	os.Mkdir(serverConfig.BanFileDir, 0777)
	for _, server := range serverList.ServerList {
		if server.Subscribed {
			log.Println("Saving ban list for server: " + server.ServerName)

			outbuf := new(bytes.Buffer)
			enc := json.NewEncoder(outbuf)
			enc.SetIndent("", "\t")

			err := enc.Encode(server.BanList)

			if err != nil {
				log.Println("Error encoding ban list file: " + err.Error())
				os.Exit(1)
			}
			err = ioutil.WriteFile(defaultBanFileDir+"/"+FileNameFilter(server.ServerName)+".json", outbuf.Bytes(), 0644)
			if err != nil {
				log.Println("Error saving ban list: " + err.Error())
				continue
			}
		}
	}
}

//Refresh list of servers from master
func updateServerList() {

	//Update server list
	log.Println("Updating server list")
	data, err := fetchFile(serverConfig.ListURL)
	if err != nil {
		log.Println("Error updating server list: " + err.Error())
	}

	//If there appears to be data, attempt parse it
	if len(data) > 0 {
		var sList serverListData

		//handle 404 error (TODO: check headers)
		lData := strings.ToLower(string(data))
		if strings.Contains(lData, "404: not found") {
			log.Println("Error updating server list: 404: Not Found")
			return
		}

		//Decode json
		err = json.Unmarshal([]byte(data), &sList)
		found := false
		if err == nil {
			//Check the new data against our current list
			for _, server := range sList.ServerList {
				foundl := false
				if server.ServerName != "" && server.ServerURL != "" {
					for ipos, s := range serverList.ServerList {
						//Found existing entry
						if s.ServerName == server.ServerName {
							foundl = true
							found = true
							//Update entry
							serverList.ServerList[ipos].ServerURL = server.ServerURL
							serverList.ServerList[ipos].Discord = server.Discord
							serverList.ServerList[ipos].Website = server.Website
							serverList.ServerList[ipos].Logs = server.Logs
							serverList.ServerList[ipos].JsonGzip = server.JsonGzip

						}
					}
					if !foundl {
						//New entry, subscribe automatically if chosen
						if serverConfig.AutoSubscribe {
							server.Subscribed = true
						} else {
							server.Subscribed = false
						}
						//Add, datestamp
						server.LocalAdd = time.Now().Format(timeFormat)
						serverList.ServerList = append(serverList.ServerList, server)
						log.Println("Added server: " + server.ServerName)
						writeServerListFile()
					}
				}
			}
			if !found {
				//Write updated file and update webServer caches if needed
				writeServerListFile()
			}
		} else {
			log.Println("Unable to parse remote server list file")
		}
	}
}

//Download file to byte array
func fetchFile(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	output, err := ioutil.ReadAll(resp.Body)

	return output, err
}

//Monitor ban file for changes
func watchBanFile() {
	var err error

	filePath := serverConfig.BanFile
	//Save current profile
	if initialStat == nil {
		initialStat, err = os.Stat(filePath)
	}

	if err != nil {
		log.Println("watchBanFile: stat: " + err.Error())
		return
	}

	if initialStat != nil {
		stat, errb := os.Stat(filePath)
		if errb != nil {
			log.Println("watchDatabaseFile: restat")
			return
		}

		//Detect file change
		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			log.Println("watchBanFile: file changed")
			readServerBanList() //Reload ban list

			//Update stat for next time
			initialStat, err = os.Stat(filePath)

			if err != nil {
				log.Println("watchBanFile: stat: " + err.Error())
				return
			}
			return
		}
	}
}
