package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

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
		foundSelf := false
		if err == nil {
			//Check the new data against our current list
			for _, server := range sList.ServerList {
				foundl := false
				if server.ServerName != "" && server.ServerURL != "" {
					for spos, s := range serverList.ServerList {
						//Found existing entry
						if s.ServerName == server.ServerName {
							foundl = true
							found = true
							//Update entry
							serverList.ServerList[spos].ServerURL = server.ServerURL
							serverList.ServerList[spos].JsonGz = server.JsonGz
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
						server.AddedLocally = time.Now()
						serverList.ServerList = append(serverList.ServerList, server)
						log.Println("Added server: " + server.ServerName)
					}
				}
				//Found ourselves in plublic list
				if server.ServerName == server.ServerName {
					foundSelf = true
				}
			}
			//Alert if we're not in the list
			if !foundSelf {
				log.Println("We are currently not found in the server list!")
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

//Incomplete
//Fetch and update a ban list from a server
func updateBanList() {
	for _, server := range serverList.ServerList {
		if server.Subscribed {
			log.Println("Updating ban list for server: " + server.ServerName)
			data, err := fetchFile(server.ServerURL)
			if err != nil {
				log.Println("Error updating ban list: " + err.Error())
			}
			if len(data) > 0 {
				//
			}
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
func WatchBanFile() {
	var err error

	filePath := serverConfig.BanFile
	//Save current profile
	if initialStat == nil {
		initialStat, err = os.Stat(filePath)
	}

	if err != nil {
		log.Println("WatchBanFile: stat: " + err.Error())
		return
	}

	if initialStat != nil {
		stat, errb := os.Stat(filePath)
		if errb != nil {
			log.Println("WatchDatabaseFile: restat")
			return
		}

		//Detect file change
		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			log.Println("WatchBanFile: file changed")
			readServerBanList() //Reload ban list

			//Update stat for next time
			initialStat, err = os.Stat(filePath)

			if err != nil {
				log.Println("WatchBanFile: stat: " + err.Error())
				return
			}
			return
		}
	}
}
