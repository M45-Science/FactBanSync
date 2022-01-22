package main

import (
	"log"
	"sort"
	"strconv"
	"time"
)

func compositeBans() {
	var compositeBanlist []banDataType

	//Add our bans first
	for _, ban := range ourBanData {

		compositeBanlist = append(compositeBanlist, banDataType{
			UserName: ban.UserName,
			Sources:  []string{serverConfig.Name},
			Reasons:  []string{ban.Reason},
			Revokes:  []bool{ban.Revoked},
			Adds:     []string{ban.Added}})
	}

	dupes := 0

	//Now add the bans from other servers
	for _, server := range serverList.ServerList {
		//Only if subscribed, and skip self
		if server.LocalData.Subscribed && server.Name != serverConfig.Name {
			for _, ban := range server.LocalData.BanList {
				//Don't composite revoked bans
				if !ban.Revoked {
					found := false
					for ipos, iban := range compositeBanlist {
						//Duplicate, add data to existing entry
						if iban.UserName == ban.UserName {
							found = true
							dupes++
							compositeBanlist[ipos].Sources = append(compositeBanlist[ipos].Sources, server.Name)

							if server.LocalData.StripReasons {
								compositeBanlist[ipos].Reasons = append(compositeBanlist[ipos].Reasons, "")
							} else {
								compositeBanlist[ipos].Reasons = append(compositeBanlist[ipos].Reasons, ban.Reason)
							}

							compositeBanlist[ipos].Revokes = append(compositeBanlist[ipos].Revokes, ban.Revoked)

							compositeBanlist[ipos].Adds = append(compositeBanlist[ipos].Adds, ban.Added)
							break
						}
					}
					if server.LocalData.StripReasons {
						ban.Reason = ""
					}
					//This isn't already in the list, so add it (avoid dupes)
					if !found {
						compositeBanlist = append(compositeBanlist, banDataType{
							UserName: ban.UserName,
							Sources:  []string{server.Name},
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
	compBan := []banDataType{}
	for bpos, ban := range compositeBanlist {
		if bpos < serverConfig.ServerPrefs.MaxBanOutputSize {
			compBan = append(compBan, ban)
		} else {
			log.Println("Banlist size (" + strconv.Itoa(serverConfig.ServerPrefs.MaxBanOutputSize) + ") exceeded, truncating...")
			break
		}
	}
	log.Println("Composite banlist updated: " + strconv.Itoa(len(compBan)) + " bans")

	var condList []minBanDataType
	//Minimize and output as a Factorio-friendly list
	for _, ban := range compBan {
		if !ban.Revoked {
			reasonList := ""
			for rpos, reason := range ban.Reasons {
				if rpos > 0 {
					reasonList += ", "
				}
				if reason != "" {
					reasonList += ban.Sources[rpos] + ": " + reason
				} else {
					reasonList += "(" + ban.Sources[rpos] + ")"
				}
			}
			condList = append(condList, minBanDataType{
				UserName: ban.UserName, Reason: reasonList})
		}
	}
	compositeBanData = condList

	writeCompositeBanlist()
}
