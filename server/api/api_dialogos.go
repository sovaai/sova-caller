package api

import (
	"path/filepath"
	"sova-caller-backend/constants"
	"sova-caller-backend/utils"
	"time"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

/* TYPES */

type DosWorkflow struct {
	CAPI                string
	UUID                string
	CUID                string
	SESSION_ID          string
	asrClient           WSClientInterface
	data                chan []byte
	asrResponse         chan string
	asrResponseReceived chan string
	needClose           chan string
	nC                  string
	error               chan string
	GREETING_MESSAGE    string
	CAMPAIGN_ID         string
	ASR                 ServiceConfig
	TTS                 ServiceConfig
}

type ChatInitParams struct {
	UUID string  `json:"uuid"`
	CUID *string `json:"cuid"`
}

type ChatReadyEventParams struct {
	CUID string `json:"cuid"`
	EUID string `json:"euid"`
}

type ChatRequestParams struct {
	CUID string `json:"cuid"`
	TEXT string `json:"text"`
}

type ChatInitType struct {
	ID     string `json:"id"`
	Result struct {
		Name   string `json:"name"`
		Cuid   string `json:"cuid"`
		Events struct {
		} `json:"events"`
	} `json:"result"`
}

type ChatResponseType struct {
	ID     string `json:"id"`
	Result struct {
		Text struct {
			Value    string `json:"value"`
			Delay    int    `json:"delay"`
			ShowRate bool   `json:"showRate"`
			Status   int    `json:"status"`
		} `json:"text"`
	} `json:"result"`
}

type NewDosWorkflowData struct {
	CAPI        string
	UUID        string
	CAMPAIGN_ID string
	SESSION_ID  string
	ASR         ServiceConfig
	TTS         ServiceConfig
}

/* TYPES */

func (dialogOS *DosWorkflow) ChatReady() string {

	const EUID = "00b2fcbe-f27f-437b-a0d5-91072d840ed3"

	params := &ChatReadyEventParams{CUID: dialogOS.CUID, EUID: EUID}
	jsonBytes, _ := json.Marshal(params)

	config := CreateRequestConfig(jsonBytes, dialogOS.CAPI, "Chat.event", "POST")

	res, _ := Request(config)

	var result ChatResponseType
	json.Unmarshal(res, &result)

	msg := result.Result.Text.Value

	return msg

}

func (dialogOS *DosWorkflow) ChatInit() string {

	defer utils.TimeTrack(time.Now(), "DIALOGOS SESSION - CHAT INIT")

	params := &ChatInitParams{UUID: dialogOS.UUID, CUID: nil}
	jsonBytes, _ := json.Marshal(params)

	config := CreateRequestConfig(jsonBytes, dialogOS.CAPI, "Chat.init", "POST")

	fmt.Println(config.URL)

	res, _ := Request(config)

	var result ChatInitType
	json.Unmarshal(res, &result)

	cuid := result.Result.Cuid

	return cuid

}

func (dialogOS *DosWorkflow) ChatRequest(text string) (string, error) {
	defer utils.TimeTrack(time.Now(), "DIALOGOS SESSION - CHAT REQUEST")

	params := &ChatRequestParams{CUID: dialogOS.CUID, TEXT: text}
	jsonBytes, _ := json.Marshal(params)

	config := CreateRequestConfig(jsonBytes, dialogOS.CAPI, "Chat.request", "POST")

	res, _ := Request(config)

	var result ChatResponseType
	json.Unmarshal(res, &result)
	msg := result.Result.Text.Value

	return msg, nil
}

func NewDosWorkflow(data *NewDosWorkflowData) *DosWorkflow {

	dialogOS := &DosWorkflow{
		UUID:        data.UUID,
		SESSION_ID:  data.SESSION_ID,
		CAPI:        data.CAPI,
		CAMPAIGN_ID: data.CAMPAIGN_ID,
		ASR:         data.ASR,
		TTS:         data.TTS,
	}

	cuid := dialogOS.ChatInit()
	dialogOS.CUID = cuid

	dialogOS.GREETING_MESSAGE = dialogOS.ChatReady()
	fmt.Printf("CHAT READY WAS SENT. RESPONSE: %s", dialogOS.GREETING_MESSAGE)

	return dialogOS
}

func (dialogOS *DosWorkflow) Start(needClose chan string, asrResponseReceived chan string, error chan string) {

	dialogOS.data = make(chan []byte)
	dialogOS.asrResponse = make(chan string)
	dialogOS.needClose = needClose
	dialogOS.asrResponseReceived = asrResponseReceived

	dialogOS.error = error

	dialogOS.asrClient = &WSAsrClient{
		URL:   dialogOS.ASR.URL,
		TOKEN: dialogOS.ASR.TOKEN,
	}

	go dialogOS.asrClient.Exec(dialogOS.data, dialogOS.asrResponse, error)

	go dialogOS.readFile()
	go dialogOS.manageCallFlow()
}

func (dialogOS *DosWorkflow) readFile() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("PANIC occurred in dosworkflow readFile method:", err)
			dialogOS.error <- "true"
		}
	}()

	asrPath := filepath.Join(constants.AUDIO_DIR, dialogOS.CAMPAIGN_ID, constants.ASR_FILE_NAME)
	file, err := os.OpenFile(asrPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println("Error opening file!")
	}

	defer file.Close()

	var offset int64 = 0
	const maxSz = 4096

	b := make([]byte, maxSz)

	for {
		res, err := file.ReadAt(b, offset)
		if err != nil && err != io.EOF {
			break
		}

		if err == io.EOF && dialogOS.nC == "need to close" {
			break
		}

		if err == nil {
			offset += int64(res)
			dialogOS.data <- b
		}
	}
}

func (dialogOS *DosWorkflow) manageCallFlow() {
	dialogOS.handleAssistantResponse(dialogOS.GREETING_MESSAGE)

	for {
		select {
		case res := <-dialogOS.asrResponse:
			fmt.Printf("Абонент:\n%s\n", res)

			answer, err := dialogOS.ChatRequest(res)
			fmt.Printf("Ассистент:\n%s\n", answer)

			if err != nil {
				dialogOS.error <- "true"
				break
			}

			dialogOS.handleAssistantResponse(answer)

		case dialogOS.nC = <-dialogOS.needClose:
			fmt.Println("dos workflow need close")

			if dialogOS.asrClient != nil {
				dialogOS.asrClient.Close()
			}

			return
		}
	}
}

func (dialogOS *DosWorkflow) createAudioFile(text string) bool {
	dec, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return false
	}

	ttsPath := filepath.Join(constants.AUDIO_DIR, dialogOS.CAMPAIGN_ID, constants.TTS_FILE_NAME)
	f, err := os.Create(ttsPath)
	if err != nil {
		return false
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return false
	}

	return true
}

func (dialogOS *DosWorkflow) handleAssistantResponse(response string) {
	audioData, err := Tts(response, dialogOS.TTS)

	if err != nil {
		dialogOS.error <- "true"
		return
	}

	if !dialogOS.createAudioFile(audioData) {
		dialogOS.error <- "true"
		return
	}

	dialogOS.asrResponseReceived <- "asr resecived"
}
