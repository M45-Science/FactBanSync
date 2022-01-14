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

func updateServerList() {

	//Update server list
	log.Println("Updating server list")
	data, err := fetchFile(serverConfig.ListURL)
	if err != nil {
		log.Println("Error updating server list: " + err.Error())
	}

	if len(data) > 0 {
		var sList serverListData

		lData := strings.ToLower(string(data))
		if strings.Contains(lData, "404: not found") {
			log.Println("Error updating server list: 404: Not Found")
			return
		}
		err = json.Unmarshal([]byte(data), &sList)
		found := false
		foundSelf := false
		if err == nil {
			for _, server := range sList.ServerList {
				foundl := false
				if server.ServerName != "" && server.ServerURL != "" {
					for _, s := range serverList.ServerList {
						if s.ServerName == server.ServerName {
							foundl = true
							found = true
						}
					}
					if !foundl {
						if serverConfig.AutoSubscribe {
							server.Subscribed = true
						} else {
							server.Subscribed = false
						}
						server.AddedLocally = time.Now()
						serverList.ServerList = append(serverList.ServerList, server)
						log.Println("Added server: " + server.ServerName)
					}
				}
				if server.ServerName == server.ServerName {
					foundSelf = true
				}
			}
			if !foundSelf {
				log.Println("We are currently not found in the server list!")
			}
			if !found {
				writeServerListFile()
			}
		} else {
			log.Println("Unable to parse remote server list file")
		}
	}
}

func fetchFile(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	output, err := ioutil.ReadAll(resp.Body)

	return output, err
}

func WatchBanFile() {
	var err error

	filePath := serverConfig.BanFile
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

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			log.Println("WatchBanFile: file changed")
			readServerBanList()

			initialStat, err = os.Stat(filePath)

			if err != nil {
				log.Println("WatchBanFile: stat: " + err.Error())
				return
			}
			return
		}
	}
}
