package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Duplicate code, not great
func fetchBanLists() {
	gDirty := 0
	for spos, server := range serverList.ServerList {
		if strings.EqualFold(server.CommunityName, serverConfig.CommunityName) {
			continue
		}
		lDirty := 0
		revoked := 0
		if server.LocalData.Subscribed {
			oldList := server.LocalData.BanList

			data, err := fetchFile(server.BanListURL)
			if err != nil {
				log.Println("Error updating ban list: " + err.Error())
				continue
			}
			if len(data) > 0 {

				var names []string
				if strings.EqualFold(server.CommunityName, "RedMew") {
					count := 0
					var redMewNames []string
					if server.UseRedScrape {
						if serverConfig.ServerPrefs.VerboseLogging {
							log.Println("Scraping RedMew.")
						}
						redMewNames = GetRedMew(server.BanListURL)
					}

					if redMewNames != nil {
						for _, red := range redMewNames {
							rLen := len(red)
							if rLen > 0 && rLen < 128 {
								names = append(names, strings.ToLower(red))
								count++
							}
						}
						if serverConfig.ServerPrefs.VerboseLogging {
							log.Printf("Redmew: %v names scraped.\n", count)
						}
					} else if strings.EqualFold(server.CommunityName, "Comfy") {
						count := 0
						var comfyNames []string
						if server.UseComfyScrape {
							if serverConfig.ServerPrefs.VerboseLogging {
								log.Println("Scraping Comfy.")
							}
							comfyNames = GetComfy(server.BanListURL)
						}

						if comfyNames != nil {
							for _, comfy := range comfyNames {
								rLen := len(comfy)
								if rLen > 0 && rLen < 128 {
									names = append(names, strings.ToLower(comfy))
									count++
								}
							}
							if serverConfig.ServerPrefs.VerboseLogging {
								log.Printf("Comfy: %v names scraped.\n", count)
							}
						}
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
									if strings.EqualFold(item.UserName, name) {
										if item.Revoked {
											if serverConfig.ServerPrefs.VerboseLogging {
												log.Println(server.CommunityName + ": Revoked ban was reinstated: " + item.UserName)
											}

											serverList.ServerList[spos].LocalData.BanList[ipos].Revoked = false
											serverList.ServerList[spos].LocalData.BanList[ipos].Added = time.Now().Format(time.RFC3339)
										}
										found = true
									}
								}
								if !found {
									gDirty++
									lDirty++
									serverList.ServerList[spos].LocalData.BanList = append(serverList.ServerList[spos].LocalData.BanList, banDataType{UserName: strings.ToLower(name), Added: time.Now().Format(timeFormat)})
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
								if strings.EqualFold(ban.UserName, item.UserName) && !item.Revoked {
									if item.Revoked {
										if serverConfig.ServerPrefs.VerboseLogging {
											log.Println(server.CommunityName + ": Revoked ban was reinstated: " + item.UserName)
										}
										serverList.ServerList[spos].LocalData.BanList[ipos].Revoked = false
									}
									found = true
								}
							}
							if !found {
								gDirty++
								lDirty++
								serverList.ServerList[spos].LocalData.BanList = append(serverList.ServerList[spos].LocalData.BanList, banDataType{UserName: strings.ToLower(item.UserName), Reason: item.Reason, Added: time.Now().Format(timeFormat)})
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
							if strings.EqualFold(ban.UserName, item.UserName) {
								found = true
								break
							}
						}
						if !found {
							count++
							revoked++
							if oldLen > 30 && count > threshold {
								if serverConfig.ServerPrefs.VerboseLogging {
									log.Println("More than 10% of bans were revoked. Ban list was probably cleared, silencing printout.")
								}
								oldLen = 0
							}
							if oldLen > 0 {
								if serverConfig.ServerPrefs.VerboseLogging {
									log.Println(server.CommunityName + ": Ban for " + item.UserName + " was revoked")
								}
							}
							serverList.ServerList[spos].LocalData.BanList[ipos].Revoked = true
						}
					}
				}

			}
			if lDirty > 0 {
				log.Printf("Found %v new bans for %v\n", lDirty, server.CommunityName)
			}
			if revoked > 0 {
				log.Printf("Found %v revoked bans for %v\n", revoked, server.CommunityName)
			}

		}
	}
	if gDirty > 0 {
		saveBanLists()
	}
}

// Refresh list of servers from master
func updateServerList() {

	//Update server list
	if serverConfig.ServerPrefs.VerboseLogging {
		log.Println("Updating server list")
	}
	data, err := fetchFile(serverConfig.ServerListURL)
	if err != nil {
		log.Println("Error updating server list: " + err.Error())
	}

	//If there appears to be data, attempt parse it
	if len(data) > 0 {
		var sList serverListData

		//Decode json
		err = json.Unmarshal([]byte(data), &sList)
		if err == nil {
			updated := false
			//Check the new data against our current list
			for _, server := range sList.ServerList {
				foundl := false
				if server.CommunityName != "" && server.BanListURL != "" {
					for ipos, s := range serverList.ServerList {
						//Found existing entry
						if strings.EqualFold(s.CommunityName, server.CommunityName) {
							foundl = true

							if serverList.ServerList[ipos].BanListURL != server.BanListURL {
								serverList.ServerList[ipos].BanListURL = server.BanListURL
								updated = true
							}
							if serverList.ServerList[ipos].WhiteListURL != server.WhiteListURL {
								serverList.ServerList[ipos].WhiteListURL = server.WhiteListURL
								updated = true
							}
							if serverList.ServerList[ipos].LogFileURL != server.LogFileURL {
								serverList.ServerList[ipos].LogFileURL = server.LogFileURL
								updated = true
							}
							if serverList.ServerList[ipos].WebsiteURL != server.WebsiteURL {
								serverList.ServerList[ipos].WebsiteURL = server.WebsiteURL
								updated = true
							}
							if serverList.ServerList[ipos].DiscordURL != server.DiscordURL {
								serverList.ServerList[ipos].DiscordURL = server.DiscordURL
								updated = true
							}
							if serverList.ServerList[ipos].JsonGzip != server.JsonGzip {
								serverList.ServerList[ipos].JsonGzip = server.JsonGzip
								updated = true
							}
							if serverList.ServerList[ipos].UseRedScrape != server.UseRedScrape {
								serverList.ServerList[ipos].UseRedScrape = server.UseRedScrape
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
						log.Println("Found new community: " + server.CommunityName)
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

// Download file to byte array
func fetchFile(url string) ([]byte, error) {

	timeout := (time.Duration(serverConfig.ServerPrefs.DownloadTimeoutSeconds) * time.Second)
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

	maxSize := int64(serverConfig.ServerPrefs.DownloadSizeLimitKB * 1024)
	if serverConfig.ServerPrefs.DownloadSizeLimitKB > 0 {
		maxSize = serverConfig.ServerPrefs.DownloadSizeLimitKB * 1024
	}
	if resp.ContentLength > maxSize {
		return []byte{}, errors.New("file too large")
	}

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}

	dlSize := len(output)
	if serverConfig.ServerPrefs.VerboseLogging {
		log.Printf("Fetched file: %s (%d kb)\n", url, dlSize/1024)
	}

	return output, err
}

// Monitor ban file for changes
func watchBanFile() {
	var err error

	filePath := serverConfig.PathData.FactorioBanFile
	//Save current profile
	if initialStat == nil {
		initialStat, err = os.Stat(filePath)
	}

	if err != nil {
		log.Println("watchBanFile: stat: " + err.Error())
		initialStat = nil
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
			if serverConfig.ServerPrefs.VerboseLogging {
				log.Println("watchBanFile: file changed")
			}
			readServerBanList() //Reload ban list
			compositeBans()     //Update composite ban list
			updateWebCache()    //Update web cache

			//Update stat for next time
			initialStat, err = os.Stat(filePath)

			if err != nil {
				log.Println("watchBanFile: stat: " + err.Error())
				initialStat = nil
				return
			}
			return
		}
	}
}
