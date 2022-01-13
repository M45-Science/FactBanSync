package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func updateServerList() {
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
		err = json.Unmarshal([]byte(data), &sList)
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
				}
			}
		} else {
			log.Println("Unable to parse remote server list file")
			os.Remove(serverConfig.ServerListFile + ".tmp")
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
