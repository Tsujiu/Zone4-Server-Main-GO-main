package config

import (
	"encoding/json"
	"os"
)

type Channel struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	RunCmd string `json:"run_cmd"`
	Port  int    `json:"port"`
}

var channels []Channel

func LoadChannels(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var arr []Channel
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	channels = arr
	return nil
}

func GetChannels() []Channel {
	return channels
}

func GetChannelByID(id string) *Channel {
	for i := range channels {
		if channels[i].ID == id {
			return &channels[i]
		}
	}
	return nil
}
