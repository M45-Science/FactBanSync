package main

import (
	"io"
	"log"
	"os"
	"time"
)

///Config and start a new log file
func startLog() {
	//Make log dir
	os.Mkdir(serverConfig.PathData.LogDir+"/", 0777)

	//Open log file
	logName := time.Now().Format("2006-01-02") + ".log"
	logDesc, err := os.OpenFile(serverConfig.PathData.LogDir+"/"+logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Println("Couldn't open log file!")
	}

	mw := io.MultiWriter(os.Stdout, logDesc) //To log and stdout
	log.SetOutput(mw)

	log.SetFlags(log.Lmicroseconds | log.Lshortfile) //Show source file and line number
}
