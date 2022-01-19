# FactBanSync<br>
[![License: Unlicense](https://img.shields.io/badge/license-Unlicense-blue.svg)](http://unlicense.org/)
<br>
[![Go](https://github.com/Distortions81/FactBanSync/actions/workflows/go.yml/badge.svg)](https://github.com/Distortions81/FactBanSync/actions/workflows/go.yml)
[![ReportCard](https://github.com/Distortions81/FactBanSync/actions/workflows/report.yml/badge.svg)](https://github.com/Distortions81/FactBanSync/actions/workflows/report.yml)
[![CodeQL](https://github.com/Distortions81/FactBanSync/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/Distortions81/FactBanSync/actions/workflows/codeql-analysis.yml)
[![BinaryBuild](https://github.com/Distortions81/FactBanSync/actions/workflows/build-linux64.yml/badge.svg)](https://github.com/Distortions81/FactBanSync/actions/workflows/build-linux64.yml)<br><br>
Grabs data in a simple and secure, decentralized fashion.<br>
*This is free and unencumbered software released into the public domain.*<br>
<br>
### Compile and setup steps<br>
https://github.com/Distortions81/FactBanSync/releases<br>
<br>
Download binary, OR:<br>
1: Install GO 1.17.x: https://go.dev/dl/<br>
2: Go to the FactBanSync directory, run 'go get'<br>
3: Run 'go build', then run the FactBanSync binary.<br>

#### Setup:<br>
1: Use the setup wizard 
<br>*(or let it generate a default config, then edit the config file)*<br>
(optional) 2: Add your server to the list:<br>
https://github.com/Distortions81/Factorio-Community-List/<br>
<br>
### What currently works:<br>
Fetching list of servers<br>
Fetching bans from other servers, detecting new bans<br>
Limit output ban list size, keeping newest.<br>
Detecting when a ban is revoked.<br>
Webserver, with cached json and json.gz<br>
Reasonable download time/size limitations<br>
<br>
### What is still WIP?
Setup Wizard<br>
RCON banning live<br>
Logfile monitoring for logins<br>
Whitelists<br>
Unit tests<br>
<br>
*(ChatWire currently handles rcon/log monitoring for me.)*
