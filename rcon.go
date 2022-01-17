package main

import (
	"log"

	"github.com/Distortions81/rcon"
)

//For live disconnect/ban
func SendRCON(address string, command string, password string) {

	//Connect
	remoteConsole, err := rcon.Dial(address, password)
	if err != nil || remoteConsole == nil {
		log.Println("rcon: " + err.Error())
		return
	}

	//Write
	reqID, err := remoteConsole.Write(command)
	if err != nil {
		log.Println(err)
		return
	}
	defer remoteConsole.Close()

	//Read
	resp, respReqID, err := remoteConsole.Read()
	if err != nil {
		log.Println("rcon: " + err.Error())
		return
	}

	//Sanity check
	if reqID != respReqID {
		log.Println("Invalid response ID.")
		return
	}

	log.Println("rcon response: " + resp)
}
