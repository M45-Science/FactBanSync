package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func readServerBanList() {

	file, err := os.Open(serverConfig.BanFile)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	var bData []banDataData

	data, err := ioutil.ReadAll(file)

	var names []string
	_ = json.Unmarshal([]byte(data), &names)

	for _, name := range names {
		if name != "" {
			bData = append(bData, banDataData{UserName: name})
		}
	}

	var bans []banDataData
	_ = json.Unmarshal([]byte(data), &bans)

	for _, item := range bans {
		if item.UserName != "" {
			if item.Address == "0.0.0.0" {
				item.Address = ""
			}
			bData = append(bData, item)
		}
	}

	banData = bData

	log.Println("Read " + fmt.Sprintf("%v", len(bData)) + " bans from banlist")
}

func writeBanListFile() {
	file, err := os.Create(serverConfig.BanFile)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(banData)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	_, err = file.Write(outbuf.Bytes())

	log.Println("Wrote banlist of " + fmt.Sprintf("%v", len(banData)) + " items")
}

func writeConfigFile() {
	file, err := os.Create(configPath)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(serverConfig)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	_, err = file.Write(outbuf.Bytes())

}

func WriteServerListFile() {
	file, err := os.Create(serverConfig.ServerListFile)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	serverList.Version = "0.0.1"
	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err = enc.Encode(serverList)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	_, err = file.Write(outbuf.Bytes())
	log.Print("Wrote server list file")
}

func readServerList() {

	file, err := os.Open(serverConfig.ServerListFile)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	var sList serverListData

	data, err := ioutil.ReadAll(file)

	err = json.Unmarshal([]byte(data), &sList)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	serverList = sList
}
