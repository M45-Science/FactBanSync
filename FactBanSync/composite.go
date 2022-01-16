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

		compositeBanlist = append(compositeBanlist, banDataType{
			UserName: ban.UserName,
			Sources:  []string{serverConfig.Name},
			Reasons:  []string{ban.Reason},
			Revokes:  []bool{ban.Revoked},
			Adds:     []string{ban.Added}})
	}

	dupes := 0
	for _, server := range serverList.ServerList {
		//Only if subscribed, and skip self
		if server.Subscribed && server.Name != serverConfig.Name {
			for _, ban := range server.BanList {
				//Don't composite revoked bans
				if !ban.Revoked {
					found := false
					for ipos, iban := range compositeBanlist {
						if iban.UserName == ban.UserName {
							found = true
							dupes++
							compositeBanlist[ipos].Sources = append(compositeBanlist[ipos].Sources, server.Name)

							compositeBanlist[ipos].Reasons = append(compositeBanlist[ipos].Reasons, ban.Reason)

							compositeBanlist[ipos].Revokes = append(compositeBanlist[ipos].Revokes, ban.Revoked)

							compositeBanlist[ipos].Adds = append(compositeBanlist[ipos].Adds, ban.Added)
							break
						}
					}
					if !found {
						compositeBanlist = append(compositeBanlist, banDataType{
							UserName: ban.UserName,
							Sources:  []string{serverConfig.Name},
							Reasons:  []string{ban.Reason},
							Revokes:  []bool{ban.Revoked},
							Adds:     []string{ban.Added}})
					}
				}
			}
		}
	}

	log.Println("Composited " + strconv.Itoa(len(compositeBanlist)) + " bans. Overlap: " + strconv.Itoa(dupes))

	//Sort by time added, new to old
	sort.Slice(compositeBanlist, func(i, j int) bool {
		//Use newest date, if there are multiple sources
		var newest_a time.Time = time.Time{}
		for _, addA := range compositeBanlist[j].Adds {
			bTime, errb := time.Parse(timeFormat, addA)
			if errb == nil {
				if bTime.After(newest_a) || newest_a.IsZero() {
					newest_a = bTime
				}
			}
		}
		var newest_b time.Time = time.Time{}
		for _, addB := range compositeBanlist[j].Adds {
			bTime, errb := time.Parse(timeFormat, addB)
			if errb == nil {
				if bTime.After(newest_b) || newest_b.IsZero() {
					newest_b = bTime
				}
			}
		}
		return newest_a.Before(newest_b)
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
