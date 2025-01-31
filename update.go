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

	//Cycle all servers
	for spos, server := range serverList.ServerList {
		//Not ourself
		if strings.EqualFold(server.CommunityName, serverConfig.CommunityName) {
			continue
		}
		lDirty := 0
		revoked := 0
		scraper := false
		//Are subscribed
		if server.LocalData.Subscribed {

			//Save previous version to detect changes
			var oldList []banDataType
			for _, old := range server.LocalData.BanList {
				oldList = append(oldList, banDataType{
					UserName: old.UserName, Reason: old.Reason, Revoked: old.Revoked,
					Added: old.Added, Sources: old.Sources, Reasons: old.Reasons, Revokes: old.Revokes, Adds: old.Adds})
			}

			//Fetch URL
			data, err := fetchFile(server.BanListURL)
			if err != nil {
				log.Printf("Error updating ban list: %v: %v: %v\n", server.CommunityName, server.BanListURL, err.Error())
				continue
			}
			if len(data) <= 0 {
				continue
			}

			//Web scrapers
			var names []string
			if server.UseRedScrape {
				names = append(names, ScrapeRedMew(server, data)...)
				scraper = true
			} else if server.UseComfyScrape {
				names = append(names, ScrapeComfy(server, data)...)
				scraper = true
			} else { //Handle array format bans
				//json.Unmarshal(data, &names)
			}

			//Process names
			if len(names) > 0 {
				//Look for new names
				for _, name := range names {
					found := false
					if name != "" {
						for i, item := range oldList {
							if strings.EqualFold(item.UserName, name) {
								if item.Revoked {
									if serverConfig.ServerPrefs.VerboseLogging {
										log.Println(server.CommunityName + ": Revoked ban was reinstated: " + item.UserName)
									}
									serverList.ServerList[spos].LocalData.BanList[i].Revoked = false
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
				//Look for names that disappeared
				for i, item := range oldList {
					found := false
					for _, name := range names {
						if strings.EqualFold(item.UserName, name) {
							found = true
							break
						}
					}
					if !found && !item.Revoked {
						gDirty++
						lDirty++
						serverList.ServerList[spos].LocalData.BanList[i].Revoked = true
						if serverConfig.ServerPrefs.VerboseLogging {
							log.Println(server.CommunityName + ": Ban for " + item.UserName + " was revoked")
							revoked++
						}
					}
				}
			}

			//Not scraper or array format ban
			if !scraper {

				var bans []banDataType
				err = json.Unmarshal(data, &bans)

				if err != nil {
					fmt.Print("") //Remove annoying warning
				}

				//Deal with nested json
				if len(bans) <= 1 && len(data) > 0 {
					fmt.Printf("Fixing nested JSON for %v.\n", server.CommunityName)
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
						for _, ban := range oldList {
							if strings.EqualFold(ban.UserName, item.UserName) {
								if item.Revoked {
									if serverConfig.ServerPrefs.VerboseLogging {
										log.Println(server.CommunityName + ": Revoked ban was reinstated: " + item.UserName)
									}
									serverList.ServerList[spos].LocalData.BanList[ipos].Revoked = false
									revoked++
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
				for ipos, item := range oldList {
					found := false
					for _, ban := range bans {
						if strings.EqualFold(ban.UserName, item.UserName) {
							found = true
							break
						}
					}
					if !found && !item.Revoked {
						revoked++
						if serverConfig.ServerPrefs.VerboseLogging {
							log.Println(server.CommunityName + ": Ban for " + item.UserName + " was revoked")
						}
						serverList.ServerList[spos].LocalData.BanList[ipos].Revoked = true
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
