package api

import (
	"encoding/json"
	"fmt"
	"log"

	"sova-caller-backend/constants"

	"github.com/gorilla/websocket"
)

type ServiceConfig struct {
	URL   string `json:"url"`
	TOKEN string `json:"token"`
}

type AuthConfig struct {
	AUTH_TYPE   string `json:"auth_type"`
	AUTH_TOKEN  string `json:"auth_token"`
	SAMPLE_RATE uint16 `json:"sample_rate"`
}

type Alternative struct {
	Text       string  `json:"text"`
	Confidence float32 `json:"confidence"`
}

type WsAsrResponseResult struct {
	Event        int    `json:"event"`
	Final        bool   `json:"final"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Alternatives []struct {
		Text       string  `json:"text"`
		Confidence float32 `json:"confidence"`
	}
}

type WsAsrResponse struct {
	Response int                   `json:"response"`
	Results  []WsAsrResponseResult `json:"results"`
}

type WSClientInterface interface {
	connect() bool
	receive()
	close()
	handleMessage(m []byte) string
	Close()
	Send(m []byte)
	Exec(c chan []byte, response chan string, error chan string)
}
type WSAsrClient struct {
	error  chan string
	conn   *websocket.Conn
	finish chan interface{}
	text   chan string
	URL    string
	TOKEN  string
}

func (w *WSAsrClient) connect() bool {
	conn, _, err := websocket.DefaultDialer.Dial(w.URL, nil)
	if err != nil {
		log.Println("Error connecting to ASR Websocket Server:", err)
		w.error <- "true"
		return false
	}
	w.conn = conn

	authConfig := AuthConfig{
		AUTH_TYPE:   "Bearer",
		AUTH_TOKEN:  w.TOKEN,
		SAMPLE_RATE: constants.DEFAULT_SAMPLE_RATE,
	}

	config, err := json.Marshal(authConfig)
	if err != nil {
		log.Println("Error Marshal auth config ASR Websocket Server:", err)
		w.error <- "true"
		return false
	}

	err = conn.WriteMessage(websocket.TextMessage, config)
	if err != nil {
		log.Println("Error auth:", err)
		w.error <- "true"
		return false
	}

	return true
}

func (w *WSAsrClient) close() {
	log.Println("close ws connection")
	if w.conn != nil {
		w.conn.Close()
	}
}

func (w *WSAsrClient) Close() {
	w.finish <- "fihish"
}

func (w *WSAsrClient) receive() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred in AsrClient receive method:", err)
			w.error <- "true"
		}
	}()

	for {
		_, msg, err := w.conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			w.error <- fmt.Sprintf("Error in receive:%s\n", err)
			return
		}
		text := w.handleMessage(msg)
		if len(text) > 0 {
			w.text <- text
		}
	}
}

func (w *WSAsrClient) Send(m []byte) {
	if w.conn != nil {
		err := w.conn.WriteMessage(websocket.BinaryMessage, m)
		if err != nil {
			log.Println("Error in send:", err)
			w.error <- "true"
			return
		}
	}
}

func (w *WSAsrClient) handleMessage(m []byte) string {
	var res WsAsrResponse
	err := json.Unmarshal(m, &res)

	if err != nil {
		panic("Message unmarshaling failed")
	}

	if len(res.Results) > 0 {
		final := findFinalResult(res.Results)
		if len(final) > 0 {
			fmt.Println(final)
			return findMaxConfidenceText(final)
		}
	}

	return ""
}

func findFinalResult(res []WsAsrResponseResult) []WsAsrResponseResult {
	var finalResults []WsAsrResponseResult
	for _, r := range res {

		if isFinalResult(r) {
			finalResults = append(finalResults, r)
		}
	}

	return finalResults
}

func isFinalResult(r WsAsrResponseResult) bool {
	return r.Final
}

func findMaxConfidenceText(res []WsAsrResponseResult) string {
	var alternatives = res[len(res)-1].Alternatives
	var maxAlt = alternatives[0]

	for _, a := range alternatives {
		if a.Confidence > maxAlt.Confidence {
			maxAlt = a
		}
	}

	return maxAlt.Text
}

func (w *WSAsrClient) Exec(c chan []byte, response chan string, error chan string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred in AsrClient Exec method:", err)
			w.error <- "true"
		}
	}()

	w.error = error
	w.finish = make(chan interface{})
	w.text = response

	if w.connect() {
		go w.receive()
	}

	for {
		select {
		case <-w.finish:
			log.Println("Finish ASR ws connection")
			return
		case data := <-c:
			w.Send(data)
		}
	}
}
