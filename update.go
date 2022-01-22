package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//Duplicate code, not great
func fetchBanLists() {
	gDirty := 0
	for spos, server := range serverList.ServerList {
		if server.Name == serverConfig.Name {
			continue
		}
		lDirty := 0
		revoked := 0
		if server.LocalData.Subscribed {
			oldList := server.LocalData.BanList

			data, err := fetchFile(server.Bans)
			if err != nil {
				log.Println("Error updating ban list: " + err.Error())
				continue
			}
			if len(data) > 0 {

				var names []string
				if server.Name == "RedMew" {
					count := 0
					var redMewNames []string
					if server.UseRedScrape {
						log.Println("Scraping RedMew.")
						redMewNames = GetRedMew(server.Bans)
					}

					if redMewNames != nil {
						for _, red := range redMewNames {
							rLen := len(red)
							if rLen > 0 && rLen < 48 {
								names = append(names, red)
								count++
							}
						}
						fmt.Println("RedMew: " + strconv.Itoa(count) + " names scraped.")
					}
				} else {
					err = json.Unmarshal(data, &names)
				}

				if err != nil {
					//Not really an error, just empty array
					//Only needed because Factorio will write some bans as an array for some unknown reason.
				} else {

					//Read bans in array format
					found := false
					if len(names) > 0 {
						for _, name := range names {
							if name != "" {
								for ipos, item := range server.LocalData.BanList {
									if item.UserName == name {
										if item.Revoked {
											log.Println(server.Name + ": Revoked ban was reinstated: " + item.UserName)

											serverList.ServerList[spos].LocalData.BanList[ipos].Revoked = false
											serverList.ServerList[spos].LocalData.BanList[ipos].Added = time.Now().Format(time.RFC3339)
										}
										found = true
									}
								}
								if !found {
									gDirty++
									lDirty++
									serverList.ServerList[spos].LocalData.BanList = append(serverList.ServerList[spos].LocalData.BanList, banDataType{UserName: name, Added: time.Now().Format(timeFormat)})
								}
							}
						}
					}
				}

				if !server.UseRedScrape {

					var bans []banDataType
					err = json.Unmarshal(data, &bans)

					if err != nil {
						fmt.Print("") //Remove annoying warning
					}

					//Deal with nested json
					if len(bans) <= 1 && len(data) > 0 {
						datastr := string(data)
						datastr = strings.TrimSuffix(datastr, "]")
						datastr = strings.TrimPrefix(datastr, "[")

						err = json.Unmarshal([]byte(datastr), &bans)
						if err != nil {
							log.Println("Error parsing ban list: " + err.Error())
						}
					}

					//Read bans in standard format
					for ipos, item := range bans {
						if item.UserName != "" {
							found := false
							for _, ban := range server.LocalData.BanList {
								if ban.UserName == item.UserName && !item.Revoked {
									if item.Revoked {
										log.Println(server.Name + ": Revoked ban was reinstated: " + item.UserName)
										serverList.ServerList[spos].LocalData.BanList[ipos].Revoked = false
									}
									found = true
								}
							}
							if !found {
								gDirty++
								lDirty++
								serverList.ServerList[spos].LocalData.BanList = append(serverList.ServerList[spos].LocalData.BanList, banDataType{UserName: item.UserName, Reason: item.Reason, Added: time.Now().Format(timeFormat)})
							}
						}
					}

					//Detect bans that were revoked
					oldLen := len(oldList)
					threshold := oldLen / 10

					count := 0
					for ipos, item := range oldList {
						found := false
						for _, ban := range server.LocalData.BanList {
							if ban.UserName == item.UserName {
								found = true
								break
							}
						}
						if !found {
							count++
							revoked++
							if oldLen > 30 && count > threshold {
								log.Println("More than 10% of bans were revoked. Ban list was probably cleared, silencing printout.")
								oldLen = 0
							}
							if oldLen > 0 {
								log.Println(server.Name + ": Ban for " + item.UserName + " was revoked")
							}
							serverList.ServerList[spos].LocalData.BanList[ipos].Revoked = true
						}
					}
				}

			}
			if lDirty > 0 {
				log.Printf("Found %v new bans for %v\n", lDirty, server.Name)
			}
			if revoked > 0 {
				log.Printf("Found %v revoked bans for %v\n", revoked, server.Name)
			}

		}
	}
	if gDirty > 0 {
		saveBanLists()
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
		if err == nil {
			updated := false
			//Check the new data against our current list
			for _, server := range sList.ServerList {
				foundl := false
				if server.Name != "" && server.Bans != "" {
					for ipos, s := range serverList.ServerList {
						//Found existing entry
						if s.Name == server.Name {
							foundl = true

							if serverList.ServerList[ipos].Bans != server.Bans {
								serverList.ServerList[ipos].Bans = server.Bans
								updated = true
							}
							if serverList.ServerList[ipos].Discord != server.Discord {
								serverList.ServerList[ipos].Discord = server.Discord
								updated = true
							}
							if serverList.ServerList[ipos].Website != server.Website {
								serverList.ServerList[ipos].Website = server.Website
								updated = true
							}
							if serverList.ServerList[ipos].Logs != server.Logs {
								serverList.ServerList[ipos].Logs = server.Logs
								updated = true
							}
							if serverList.ServerList[ipos].JsonGzip != server.JsonGzip {
								serverList.ServerList[ipos].JsonGzip = server.JsonGzip
								updated = true
							}

						}
					}
					if !foundl {
						//New entry, subscribe automatically if chosen
						if serverConfig.ServerPrefs.AutoSubscribe {
							server.LocalData.Subscribed = true
						} else {
							server.LocalData.Subscribed = false
						}
						//Add, datestamp
						server.LocalData.Added = time.Now().Format(timeFormat)
						serverList.ServerList = append(serverList.ServerList, server)
						updated = true
						log.Println("Added server: " + server.Name)
						writeServerListFile()
					}
				}
			}
			if updated {
				writeServerListFile()
			}
		} else {
			log.Println("Unable to parse remote server list file")
		}
	}
}

//Download file to byte array
func fetchFile(url string) ([]byte, error) {

	timeout := defaultDownloadTimeout * time.Second
	if serverConfig.ServerPrefs.DownloadTimeoutSeconds > 0 {
		timeout = (time.Duration(serverConfig.ServerPrefs.DownloadTimeoutSeconds) * time.Second)
	}
	c := &http.Client{
		Timeout: timeout,
	}

	resp, err := c.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []byte{}, errors.New("HTTP Error: " + resp.Status)
	}

	maxSize := int64(defaultMaxDownloadSize)
	if serverConfig.ServerPrefs.DownloadSizeLimitBytes > 0 {
		maxSize = serverConfig.ServerPrefs.DownloadSizeLimitBytes
	}
	if resp.ContentLength > maxSize {
		return []byte{}, errors.New("file too large")
	}

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}

	dlSize := len(output)
	log.Println("Fetched file: " + url + " (" + strconv.Itoa(dlSize/1024) + " kb)")

	return output, err
}

//Monitor ban file for changes
func watchBanFile() {
	var err error

	filePath := serverConfig.PathData.FactorioBanFile
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
			updateWebCache()    //Update web cache
			compositeBans()     //Update composite ban list

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
