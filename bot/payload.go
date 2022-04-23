package bot

import (
	"encoding/json"

	"github.com/Toffee-iZt/HwBot/vkapi"
)

func unwrap(payload vkapi.JSONData) (string, vkapi.JSONData) {
	var botPayload struct {
		HwBot string         `json:"hwbot"`
		Data  vkapi.JSONData `json:"data"`
	}
	err := json.Unmarshal([]byte(payload), &botPayload)
	if err != nil {
		return "", ""
	}
	return botPayload.HwBot, botPayload.Data
}

func wrap(mod *Module, payload interface{}) vkapi.JSONData {
	d, _ := vkapi.NewJSONData(struct {
		HwBot string      `json:"hwbot"`
		Data  interface{} `json:"data"`
	}{
		HwBot: mod.Name,
		Data:  payload,
	})
	return d
}
