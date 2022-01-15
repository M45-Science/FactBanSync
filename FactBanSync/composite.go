package main

import (
	"log"
	"sort"
	"strconv"
)

func compositeBans() {
	var compositeBanlist []banDataType
	for _, ban := range ourBanData {
		compositeBanlist = append(compositeBanlist, ban)
	}

	dupes := 0
	for _, server := range serverList.ServerList {
		if server.Subscribed && server.Name != serverConfig.Name {
			for _, ban := range server.BanList {
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

	log.Println("Found " + strconv.Itoa(len(compositeBanlist)) + " bans, " + strconv.Itoa(dupes) + " duplicates")

	//Sort by local add
	sort.Slice(compositeBanlist, func(i, j int) bool {
		return compositeBanlist[i].LocalAdd > compositeBanlist[j].LocalAdd
	})

	compBanData = []banDataType{}
	for bpos, ban := range compositeBanlist {
		if bpos < serverConfig.MaxBanlistSize {
			compBanData = append(compBanData, ban)
		} else {
			log.Println("Banlist size exceeded, truncating")
			break
		}
	}
	log.Println("Composite banlist updated: " + strconv.Itoa(len(compBanData)) + " bans")
}
