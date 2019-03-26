package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// getLocalMedia reads media from local directory
func getLocalMedia() (media []Media) {
	files, err := ioutil.ReadDir("./media")
	if err != nil {
		log.Fatal("Error reading directory", err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".jpg") {
			m := Media{}
			m.LocalFile = file.Name()
			media = append(media, m)
		}
	}
	return media
}

// getRemoteMedia reads media from json file
func getRemoteMedia() (media []Media) {
	// check if file exists, return empty
	// likely scenario would be first run
	if _, err := os.Stat("media.json"); os.IsNotExist(err) {
		fmt.Println("INFO: media.json does not exist, first run?")
		return media
	}

	file, err := ioutil.ReadFile("media.json")
	if err != nil {
		fmt.Println("Error reading media.json, permissions?", err)
	} else {
		if err := json.Unmarshal(file, &media); err != nil {
			fmt.Println("Error parsing JSON from media.json", err)
		}
	}
	return media
}

func compareMedia(local, remote []Media) (media []Media) {
	for _, m := range local {
		exists := false
		for _, r := range remote {
			if m.LocalFile == r.LocalFile {
				exists = true
				fmt.Println("Skipping ", m.LocalFile)
			}
		}
		if !exists {
			media = append(media, m)
		}
	}
	return media
}

func uploadMediaItems(media []Media) []Media {
	for i, m := range media {
		m = uploadMedia(m)
		media[i].Id = m.Id
		media[i].Link = m.Link
	}
	return media
}

// writeRemoteMedia
func writeRemoteMedia(media []Media) {
	if len(media) == 0 {
		fmt.Println("No new media to write.")
		return
	}
	// append new post json
	existingMedia := getRemoteMedia()
	existingMedia = append(existingMedia, media...)

	// write file
	json, err := json.Marshal(media)
	if err != nil {
		fmt.Println("JSON Encoding Error", err)
	} else {
		err = ioutil.WriteFile("media.json", json, 0644)
		if err != nil {
			fmt.Println("Error writing media.json", err)
		} else {
			fmt.Println("media.json written")
		}
	}
}