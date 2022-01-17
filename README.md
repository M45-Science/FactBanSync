# FactBanSync<br>
[![License: Unlicense](https://img.shields.io/badge/license-Unlicense-blue.svg)](http://unlicense.org/)
[![Go](https://github.com/Distortions81/FactBanSync/actions/workflows/go.yml/badge.svg)](https://github.com/Distortions81/FactBanSync/actions/workflows/go.yml)
[![ReportCard](https://github.com/Distortions81/FactBanSync/actions/workflows/report.yml/badge.svg)](https://github.com/Distortions81/FactBanSync/actions/workflows/report.yml)

`This is currently a unfinished prototype, WIP.`<br>
*This is free and unencumbered software released into the public domain.*<br>
<br>
**(no binary releases yet)**<br>
<br>
## Compile and setup steps<br>
1: Install GO 1.17.x: https://go.dev/dl/<br>
2: Go to the FactBanSync directory, run 'go get'<br>
3: Run 'go build', then run the FactBanSync binary.<br>
4: Use the setup wizard  *(or let it generate a default config, then edit the config file)*<br>
5: (optional) Add your server to the list:<br>
https://github.com/Distortions81/Factorio-Community-List/<br>
<br>
### What currently works:<br>
Fetching list of servers<br>
Fetching bans from other servers, detecting new bans<br>
Limit output ban list size, keeping newest.<br>
Detecting when a ban is revoked.<br>
Webserver, with cached json and json.gz<br>
<br>
### What is still WIP?
Writing out composited bans.<br>
RCON banning live<br>
Logfile monitoring for logins<br>
Whitelists<br>
Unit tests.<br>
