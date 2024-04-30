package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	ClientID         string `json:"client_id"`
	ClientSecret     string `json:"client_secret"`
	DownloadFolder   string `json:"download_folder"`
	PreferredQuality string `json:"quality"`
	Streamer         string `json:"streamer"`
}

func (config *Config) readConfig() error {

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed to open config.json:", err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		log.Fatal("Failed to decode config.json:", err)
		return err
	}

	return nil
}

func main() {
	const interval = 15
	var isRecording = false

	config := Config{}
	recorder := NewRecorder()

	if err := config.readConfig(); err != nil {
		log.Fatal("Error reading config:", err)
	}

	downloadFolder, err := filepath.Abs(config.DownloadFolder)
	if err != nil {
		log.Fatal("Error getting absolute path for download folder:", err)
	}

	twitchClient, err := NewTwitch(config.ClientID, config.ClientSecret)
	if err != nil {
		log.Fatal("Failed to create Twitch client:", err)
	}

	twitchClient.PrintClientInfo()

	for {
		var isLive = twitchClient.IsLive(config.Streamer)

		if isLive && !isRecording {
			log.Printf("%s is now live!", config.Streamer)
			recorder.StartRecording(config.Streamer, downloadFolder, config.PreferredQuality)
			isRecording = true
		} else if !isLive && isRecording {
			log.Printf("%s has gone offline!", config.Streamer)
			recorder.StopRecording()
			isRecording = false
		}

		time.Sleep(interval * time.Second)
	}
}
