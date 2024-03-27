package main

import (
	"log"
	"strings"
)

func GetRedMew(data []byte) []string {

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
