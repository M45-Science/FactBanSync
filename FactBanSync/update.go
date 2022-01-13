package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func updateServerList() {
	defer os.Remove(serverConfig.ServerListFile + ".tmp")
	//Update server list
	log.Println("Updating server list")
	var wrote int64
	wrote, err := downloadFile(serverConfig.ServerListFile+".tmp", defaultListURL)
	if err != nil {
		log.Println("Error updating server list: " + err.Error())
	}

	if wrote > 0 {
		var sList serverListData
		data, err := ioutil.ReadFile(serverConfig.ServerListFile + ".tmp")
		if err != nil {
			log.Println("Error reading server list: " + err.Error())
		}

		lData := strings.ToLower(string(data))
		if strings.Contains(lData, "404: not found") {
			log.Println("Error updating server list: 404: Not Found")
			return
		}
		err = json.Unmarshal([]byte(data), &sList)
		found := false
		if err == nil {
			for _, server := range sList.ServerList {
				if server.ServerName != "" && server.ServerURL != "" {
					if serverConfig.AutoSubscribe {
						server.Subscribed = true
					} else {
						server.Subscribed = false
					}
					server.Added = time.Now()
					serverList.ServerList = append(serverList.ServerList, server)
					log.Println("Added server: " + server.ServerName)
					found = true
				}
			}
			if found {
				writeServerListFile()
				writeBanListFile()
			}
		} else {
			log.Println("Unable to parse remote server list file")
		}
	}
}

func downloadFile(filepath string, url string) (int64, error) {

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	var wrote int64
	wrote, err = io.Copy(out, resp.Body)
	return wrote, err
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
