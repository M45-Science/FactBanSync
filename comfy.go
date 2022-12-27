package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func GetComfy(url string) []string {
	resp, err := http.Get(url)

	if err != nil {
		log.Println("Error:", err)
	}

	//This will eventually break, probably -- 12/2022
	if resp.StatusCode == 200 {
		if resp.Body != nil {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error:", err)
			}
			dstr := string(data)
			spltStr := strings.SplitAfter(dstr, "<ul class=\"ul-has-multiple-rows\">")
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
	} else {
		log.Println("Error:", resp.StatusCode)
	}

	return nil

}
