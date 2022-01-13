package main

import (
	"log"

	"github.com/Distortions81/rcon"
	"github.com/bwmarrin/discordgo"
)

func SendRCON(address string, command string, s *discordgo.Session) {

	remoteConsole, err := rcon.Dial(address, serverConfig.RCONPassword)
	if err != nil || remoteConsole == nil {
		log.Println("rcon: " + err.Error())
		return
	}

	reqID, err := remoteConsole.Write(command)
	if err != nil {
		log.Println(err)
		return
	}
	defer remoteConsole.Close()

	resp, respReqID, err := remoteConsole.Read()
	if err != nil {
		log.Println("rcon: " + err.Error())
		return
	}

	if reqID != respReqID {
		log.Println("Invalid response ID.")
		return
	}

	log.Println("rcon response: " + resp)
}
