package main

import (
	"log"
	"strings"
)

func ParseRedMew(data []byte) []string {

	dstr := string(data)
	spltStr := strings.SplitAfter(dstr, "<ul>")
	if len(spltStr) <= 1 {
		log.Println("GetRedMew: Data not long enough.")
		return []string{}
	}
	spltStr = strings.SplitAfter(spltStr[1], "</ul>")
	cleanStr := strings.Replace(spltStr[0], "</ul>", "", -1)

	lines := strings.Split(cleanStr, "\n")
	for lpos := range lines {
		lines[lpos] = strings.TrimSpace(lines[lpos])
		if len(lines[lpos]) < 64 {
			lines[lpos] = strings.Replace(lines[lpos], "<li>", "", -1)
			lines[lpos] = strings.Replace(lines[lpos], "</li>", "", -1)
		} else {
			lines[lpos] = ""
		}
	}

	return lines

}

func ScrapeRedMew(server serverData, data []byte) []string {
	count := 0
	var names []string
	var redMewNames []string
	if server.UseRedScrape {
		if serverConfig.ServerPrefs.VerboseLogging {
			log.Println("Scraping RedMew.")
		}
		redMewNames = ParseRedMew(data)
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
	}

	return names
}
