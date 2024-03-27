package main

import (
	"log"
	"strings"
)

func GetComfy(data []byte) []string {

	var lines []string
	dstr := string(data)
	dstr = strings.ReplaceAll(dstr, "\n", "")
	dstr = strings.ReplaceAll(dstr, "\r", "")
	spltStr := strings.SplitAfter(dstr, "</tbody>")
	if len(spltStr) <= 0 {
		log.Println("GetComfy: Not enough data.")
		return lines
	}
	spltStr = strings.SplitAfter(spltStr[0], "<tbody>")

	if len(spltStr) <= 0 {
		log.Println("GetRedMew: Data not valid.")
		return []string{}
	}
	spltStr = strings.SplitAfter(spltStr[1], "<tr>")
	for _, item := range spltStr {
		newSplit := strings.SplitAfter(item, "<td>")
		for e, element := range newSplit {
			if e%2 == 0 {
				continue
			}
			if strings.Contains(element, ":") {
				continue
			}
			name := strings.ReplaceAll(element, " ", "")
			name = strings.ReplaceAll(name, "<td>", "")
			name = strings.ReplaceAll(name, "</td>", "")
			//fmt.Printf("Item: %v\n", name)
			if len(name) > 3 {
				lines = append(lines, name)
			}
		}
	}
	return lines
}
