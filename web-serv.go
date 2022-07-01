package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func startWebserver() {
	//Run a webserver, if requested
	exit := false
	if serverConfig.WebServer.RunWebServer {
		if serverConfig.PathData.FactorioBanFile == "" {
			log.Println("No factorio banlist file specified in config file")
			exit = true
		}
		if serverConfig.WebServer.SSLCertFile == "" || serverConfig.WebServer.SSLKeyFile == "" {
			log.Println("No SSL certificate or key file specified in config file")
			exit = true
		}
		if serverConfig.WebServer.DomainName == "" {
			log.Println("No domain name specified in config file")
			exit = true
		}
		if !exit {
			http.HandleFunc("/", handleFileRequest)
			server := &http.Server{
				Addr:         serverConfig.WebServer.DomainName + ":" + strconv.Itoa(serverConfig.WebServer.SSLWebPort),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 5 * time.Second,
				TLSConfig:    &tls.Config{ServerName: serverConfig.WebServer.DomainName},
			}
			go func(sc serverConfigData, serv *http.Server) {
				err := serv.ListenAndServeTLS(sc.WebServer.SSLCertFile, sc.WebServer.SSLKeyFile)
				if err != nil {
					log.Println(err)
				}
			}(serverConfig, server)
			if serverConfig.ServerPrefs.VerboseLogging {
				log.Println("Web server started:")
				log.Printf("https://%s:%d/%s.gz\n", serverConfig.WebServer.DomainName, serverConfig.WebServer.SSLWebPort, defaultFileWebName)
				log.Printf("https://%s:%d/%s\n", serverConfig.WebServer.DomainName, serverConfig.WebServer.SSLWebPort, defaultFileWebName)
			}
		} else {
			log.Println("Web server not started.")
		}
	}
}

//Web server
func handleFileRequest(w http.ResponseWriter, r *http.Request) {

	//Limit requests per second
	delay := 1000000 / serverConfig.WebServer.MaxRequestsPerSecond
	if delay < 1 {
		delay = 1
	}
	defer time.Sleep(time.Duration(delay) * time.Microsecond)

	//Cached gzip copy
	if r.URL.Path == "/"+defaultFileWebName+".gz" {
		if cachedBanListGz == nil {
			noDataReply(w)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "gzip")
		cachedBanListLock.Lock()
		w.Write(cachedBanListGz)
		cachedBanListLock.Unlock()

		//Cached copy
	} else if r.URL.Path == "/"+defaultFileWebName {
		if cachedBanList == nil {
			noDataReply(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		cachedBanListLock.Lock()
		w.Write(cachedBanList)
		cachedBanListLock.Unlock()
	} else {
		//Not found
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
		fmt.Fprintf(w, "404: File not found\n")
	}
}

func noDataReply(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	fmt.Fprintf(w, "No ban data\n")
}

func updateWebCache() {

	var localCopy []banDataType
	for _, item := range ourBanData {
		if item.UserName != "" {
			var name, reason string

			name = item.UserName

			if !serverConfig.ServerPrefs.StripReasons {
				reason = item.Reason
			}

			localCopy = append(localCopy, banDataType{UserName: strings.ToLower(name), Reason: reason})
		}
	}

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	err := enc.Encode(localCopy)

	if err != nil {
		log.Println("Error encoding ban list file: " + err.Error())
		os.Exit(1)
	}

	//Cache a normal and gzip version of the ban list
	if serverConfig.WebServer.RunWebServer {
		cachedBanListLock.Lock()
		cachedBanList = outbuf.Bytes()
		cachedBanListGz = compressGzip(outbuf.Bytes())
		if serverConfig.ServerPrefs.VerboseLogging {
			log.Printf("Cached response: json: %v gz: %v (bytes)\n", len(cachedBanList), len(cachedBanListGz))
		}
		cachedBanListLock.Unlock()
	}
}
