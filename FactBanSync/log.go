package main

import (
	"io"
	"log"
	"os"
	"time"
)

func startLog() {
	//Make log dir
	err := os.Mkdir(serverConfig.LogDir, 0777)

	if os.IsNotExist(err) {
		log.Println("Couldnt't create log dir!")
		os.Exit(1)
	}

	//Open log file
	logName := time.Now().Format("2006-01-02") + ".log"
	logDesc, err = os.OpenFile(serverConfig.LogDir+"/"+logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Println("Couldn't open log file!")
	}

	defer logDesc.Close()
	mw := io.MultiWriter(os.Stdout, logDesc) //To log and stdout
	log.SetOutput(mw)

	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}
