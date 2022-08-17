package api

import (
	"encoding/json"
	"fmt"
	"sova-caller-backend/constants"
	"sova-caller-backend/utils"
	"time"
)

type TtsRequest struct {
	TEXT   string  `json:"text"`
	VOICE  string  `json:"voice"`
	RATE   float32 `json:"rate,omitempty"`
	PITCH  float32 `json:"pitch,omitempty"`
	VOLUME float32 `json:"volume,omitempty"`
}
type TtsResponseType struct {
	ResponseCode int `json:"response_code"`
	Response     []struct {
		ResponseAudio string `json:"response_audio"`
	} `json:"response"`
}

func Tts(msg string, config ServiceConfig) (string, error) {
	defer utils.TimeTrack(time.Now(), "TTS")

	data := TtsRequest{
		TEXT:  msg,
		VOICE: constants.DEFAULT_TTS_VOICE,
	}

	dataJSON, err := json.Marshal(data)

	if err != nil {
		return "", fmt.Errorf("marshal tts request payload failed: %s", err)
	}

	req := RequestConfig{
		URL:         config.URL,
		HTTP_METHOD: "POST",
		PAYLOAD:     dataJSON,
		TOKEN:       fmt.Sprintf("Basic %s", config.TOKEN),
	}

	body, _ := Request(req)

	var res TtsResponseType
	err = json.Unmarshal(body, &res)

	if err != nil {
		return "", fmt.Errorf("unmarshal tts result failed: %s", err)
	}

	return res.Response[0].ResponseAudio, nil
}
