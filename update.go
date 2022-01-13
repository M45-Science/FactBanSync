package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func updateServerList() {
	//Update server list
	log.Println("Updating server list")
	var wrote int64
	wrote, err := downloadFile(serverConfig.ServerListFile+".tmp", defaultListURL)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	if wrote > 0 {
		var sList serverListData
		data, err := ioutil.ReadFile(serverConfig.ServerListFile + ".tmp")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(data), &sList)
		if err == nil {
			err = os.Rename(serverConfig.ServerListFile+".tmp", serverConfig.ServerListFile)
			if err != nil {
				panic(err)
			}
			serverList = sList
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
