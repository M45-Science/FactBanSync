package main

import (
	"log"
	"sort"
	"strconv"
	"time"
)

func compositeBans() {
	var compositeBanlist []banDataType
	for _, ban := range ourBanData {
		compositeBanlist = append(compositeBanlist, ban)
	}

	dupes := 0
	for _, server := range serverList.ServerList {
		//Only if subscribed, and skip self
		if server.Subscribed && server.Name != serverConfig.Name {
			for _, ban := range server.BanList {
				//Don't composite revoked bans
				if !ban.Revoked {
					found := false
					for _, iban := range ourBanData {
						if iban.UserName == ban.UserName {
							found = true
							dupes++
							break
						}
					}
					if !found {
						compositeBanlist = append(compositeBanlist, ban)
					}
				}
			}
		}
	}

	log.Println("Found " + strconv.Itoa(len(compositeBanlist)) + " bans, " + strconv.Itoa(dupes) + " duplicates")

	//Sort by time added, new to old
	sort.Slice(compositeBanlist, func(i, j int) bool {
		aTime, erra := time.Parse(timeFormat, compositeBanlist[i].LocalAdd)
		bTime, errb := time.Parse(timeFormat, compositeBanlist[j].LocalAdd)
		if erra != nil || errb != nil {
			log.Println("Error parsing time: " + erra.Error())
			return false
		}
		return aTime.Before(bTime)
	})

	//Cut list to size, new entries are at the start
	compBanData = []banDataType{}
	for bpos, ban := range compositeBanlist {
		if bpos < serverConfig.MaxBanlistSize {
			compBanData = append(compBanData, ban)
		} else {
			log.Println("Banlist size exceeded, truncating...")
			break
		}
	}
	log.Println("Composite banlist updated: " + strconv.Itoa(len(compBanData)) + " bans")

	writeCompositeBanlist()
}
