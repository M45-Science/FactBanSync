package main

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func compositeFixture(t *testing.T) {
	t.Helper()
	oldConfig, oldServers, oldBans, oldComposite := serverConfig, serverList, ourBanData, compositeBanData
	t.Cleanup(func() {
		serverConfig, serverList, ourBanData, compositeBanData = oldConfig, oldServers, oldBans, oldComposite
	})
	serverConfig = serverConfigData{CommunityName: "Home", PathData: filePathData{CompositeBanFile: filepath.Join(t.TempDir(), "composite.json")}, ServerPrefs: serverPrefs{MaxBanOutputCount: 100}}
	serverList = serverListData{}
	ourBanData = nil
	compositeBanData = nil
}

func TestCompositeBansMergesSubscribedSources(t *testing.T) {
	compositeFixture(t)
	ourBanData = []banDataType{{UserName: "ALICE", Reason: "local"}}
	serverList.ServerList = []serverData{
		{CommunityName: "Peer", LocalData: localData{Subscribed: true, BanList: []banDataType{{UserName: "alice", Reason: "remote"}, {UserName: "Revoked", Revoked: true}}}},
		{CommunityName: "Ignored", LocalData: localData{BanList: []banDataType{{UserName: "ignored"}}}},
		{CommunityName: "home", LocalData: localData{Subscribed: true, BanList: []banDataType{{UserName: "self"}}}},
	}
	compositeBans()
	want := []minBanDataType{{UserName: "alice", Reason: "Home: local, Peer: remote"}}
	if !reflect.DeepEqual(compositeBanData, want) {
		t.Fatalf("composite = %#v, want %#v", compositeBanData, want)
	}
}

func TestCompositeBansKeepsNewestAtOutputLimit(t *testing.T) {
	compositeFixture(t)
	serverConfig.ServerPrefs.MaxBanOutputCount = 2
	base := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	ourBanData = []banDataType{
		{UserName: "oldest", Added: base.Format(timeFormat)},
		{UserName: "middle", Added: base.Add(time.Hour).Format(timeFormat)},
		{UserName: "newest", Added: base.Add(2 * time.Hour).Format(timeFormat)},
	}
	compositeBans()
	if len(compositeBanData) != 2 || compositeBanData[0].UserName != "newest" || compositeBanData[1].UserName != "middle" {
		t.Fatalf("wrong bans retained: %#v", compositeBanData)
	}
}

func TestCompositeBansUsesNewestDuplicateTimestamp(t *testing.T) {
	compositeFixture(t)
	serverConfig.ServerPrefs.MaxBanOutputCount = 1
	base := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	ourBanData = []banDataType{{UserName: "alice", Added: base.Add(time.Hour).Format(timeFormat)}, {UserName: "bob", Added: base.Format(timeFormat)}}
	serverList.ServerList = []serverData{{CommunityName: "Peer", LocalData: localData{Subscribed: true, BanList: []banDataType{{UserName: "bob", Added: base.Add(2 * time.Hour).Format(timeFormat)}}}}}
	compositeBans()
	if len(compositeBanData) != 1 || compositeBanData[0].UserName != "bob" {
		t.Fatalf("newest remote timestamp ignored: %#v", compositeBanData)
	}
}
