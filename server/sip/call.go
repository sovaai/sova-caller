package sip

import (
	"fmt"
	"path/filepath"
	"sova-caller-backend/api"
	"sova-caller-backend/constants"
	pjsua2 "sova-caller-backend/pjsua2"
	"sync"
	"time"
)

const CHECK_STREAM_INTERVAL_MSEC = 400
const SILENCE_INTERAVAL_MSEC = 6000
const maxSilenceCounterValue = SILENCE_INTERAVAL_MSEC / CHECK_STREAM_INTERVAL_MSEC
const DEFAULT_CALLING_TIMEOUT_SEC = 20

type Call struct {
	pjsua2.Call
	state                        pjsua2.Pjsip_inv_state
	mutex                        sync.Mutex
	asrResponseReceived          chan string
	needClose                    chan string
	uuid                         string
	player                       pjsua2.AudioMediaPlayer
	audMed                       pjsua2.AudioMedia
	dosWorkflow                  *api.DosWorkflow
	error                        chan string
	ticker                       *time.Ticker
	voiceActivityOnSpeechCounter int
	silenceOnSilenceCounter      int
	CALL_ID                      string
	audioRecorder                pjsua2.AudioMediaRecorder
	CAMPAIGN_ID                  string
	SESSION_ID                   string
	callingTimer                 *time.Timer
	CALLING_TIMEOUT              int
}

type AudioMediaPlayer struct {
	pjsua2.AudioMediaPlayer
	ttsPath string
}

type NewCallData struct {
	UUID        string
	CAMPAIGN_ID string
	SESSION_ID  string
}

func NewCall(data *NewCallData) *Call {

	return &Call{
		uuid:                         data.UUID,
		voiceActivityOnSpeechCounter: 0,
		silenceOnSilenceCounter:      0,
		audioRecorder:                nil,
		CAMPAIGN_ID:                  data.CAMPAIGN_ID,
		SESSION_ID:                   data.SESSION_ID,
		player:                       nil,
		callingTimer:                 nil,
		CALLING_TIMEOUT:              DEFAULT_CALLING_TIMEOUT_SEC,
	}
}

func (call *Call) OnCallState(prm pjsua2.OnCallStateParam) {
	//CheckThread()

	defer call.mutex.Unlock()

	call.mutex.Lock()

	callInfo := call.GetInfo()

	if call.state == callInfo.GetState() {
		return
	}

	call.state = callInfo.GetState()
	call.CALL_ID = callInfo.GetCallIdString()

	fmt.Printf("OnCallState=%v, RemoteUri=%v, callId=%s, lastStatusCode=%v\n", callInfo.GetStateText(), callInfo.GetRemoteUri(), call.CALL_ID, call.GetInfo().GetLastStatusCode())

	if call.dosWorkflow == nil {
		call.interrupt()
		OnCallClosed(call.CALL_ID, "error")
		return
	}

	switch call.state {
	case pjsua2.PJSIP_INV_STATE_CALLING:
		fmt.Println("OnCallState CALLING!")
		call.startCallingTimer()
		return

	case pjsua2.PJSIP_INV_STATE_CONNECTING:
		fmt.Println("OnCallState CONNECTING!")
		call.callingTimer.Stop()
		return

	case pjsua2.PJSIP_INV_STATE_CONFIRMED:
		fmt.Println("OnCallState CONFIRMED!")
		call.dosWorkflow.Start(call.needClose, call.asrResponseReceived, call.error)
		go call.manageCallFlow()
		return

	case pjsua2.PJSIP_INV_STATE_DISCONNECTED:
		fmt.Println("OnCallState DISCONNECTED!")
		call.callingTimer.Stop()

		call.Close()
		return

	default:
		fmt.Println("OnCallState UNHANDLED STATUS!")
	}
}

func (call *Call) OnCallSdpCreated(prm pjsua2.OnCallSdpCreatedParam) {
	CheckThread()

	callInfo := call.GetInfo()

	fmt.Printf("OnCallSdpCreated, SendingSdp=%v\n", prm.GetSdp().GetWholeSdp())
	fmt.Printf("OnCallSdpCreated, sipCallId=%v, role=%v, state=%v\n",
		call.CALL_ID, callInfo.GetRole(), callInfo.GetState())

}

func (call *Call) OnCallMediaState(prm pjsua2.OnCallMediaStateParam) {
	fmt.Printf("OnCallMediaState\n")

	if call.GetInfo().GetLastStatusCode() != pjsua2.PJSIP_SC_OK {
		return
	}

	aud_rec := pjsua2.NewAudioMediaRecorder()
	asrPath := filepath.Join(constants.AUDIO_DIR, call.CAMPAIGN_ID, constants.ASR_FILE_NAME)
	aud_rec.CreateRecorder(asrPath)

	aud_med := call.GetAudioMedia(-1)

	call.audMed = aud_med
	call.audMed.StartTransmit(aud_rec)

	call.audioRecorder = pjsua2.NewAudioMediaRecorder()
	recordPath := filepath.Join(constants.AUDIO_DIR, call.CAMPAIGN_ID, constants.SESSIONS_DIR, call.SESSION_ID+".wav")
	call.audioRecorder.CreateRecorder(recordPath)
	call.audioRecorder.GetPortInfo().GetFormat().SetClockRate(constants.DEFAULT_SAMPLE_RATE)
	call.audMed.StartTransmit(call.audioRecorder)
}

func (call *Call) listenAudioStream(t time.Time) {
	CheckThread()

	if call.Call.HasMedia() && call.Call.IsActive() {
		if call.audMed.GetRxLevel() > 0 {
			call.silenceOnSilenceCounter = 0

			if call.audMed.GetTxLevel() > 0 {
				call.voiceActivityOnSpeechCounter++
			} else {
				call.voiceActivityOnSpeechCounter = 0
			}
		} else {
			call.voiceActivityOnSpeechCounter = 0
		}

		if call.audMed.GetTxLevel() > 0 {
			call.silenceOnSilenceCounter = 0
		} else {
			if call.audMed.GetRxLevel() == 0 {
				call.silenceOnSilenceCounter++
			}
		}

		if call.voiceActivityOnSpeechCounter > 1 {
			call.player.StopTransmit(call.audMed)
			call.player.StopTransmit(call.audioRecorder)
		}

		if call.silenceOnSilenceCounter >= maxSilenceCounterValue {
			call.stopListenAudioSteam()
			call.interrupt()
		}
	}
}

func (call *Call) manageCallFlow() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("PANIC occurred In Call manageCallFlow method:", err)

			call.interrupt()
			OnCallClosed(call.CALL_ID, "error")
		}
	}()

	CheckThread()

	for {
		select {
		case <-call.asrResponseReceived:
			call.playFileStream()
		case <-call.needClose:
			fmt.Println("CALL SHOULD BE CLOSED")

		case <-call.error:
			fmt.Println("ERROR IN MANAGE CALL FLOW")
			call.interrupt()
			OnCallClosed(call.CALL_ID, "error")

		case t := <-call.ticker.C:
			call.listenAudioStream(t)

		}

	}
}

func CreateMediaPlayer(ttsPath string) *AudioMediaPlayer {
	player := pjsua2.NewAudioMediaPlayer()
	return &AudioMediaPlayer{player, ttsPath}
}

func (call *Call) playFileStream() {
	CheckThread()

	if call.Call.HasMedia() && call.Call.IsActive() {
		if call.audMed.GetRxLevel() == 0 {

			if call.player != nil {
				call.player.StopTransmit(call.audMed)
				call.player.StopTransmit(call.audioRecorder)
				call.player = nil
			}

			ttsPath := filepath.Join(constants.AUDIO_DIR, call.CAMPAIGN_ID, constants.TTS_FILE_NAME)

			audioPlayer := CreateMediaPlayer(ttsPath)
			call.player = audioPlayer

			audioPlayer.CreatePlayer(ttsPath, uint(pjsua2.PJMEDIA_FILE_NO_LOOP))
			audioPlayer.StartTransmit(call.audMed)
			audioPlayer.StartTransmit(call.audioRecorder)
		}
	}
}

func (call *Call) interrupt() {
	if call.state == pjsua2.PJSIP_INV_STATE_DISCONNECTED {
		return
	}

	fmt.Println("interrupt")
	CheckThread()
	callOpParam := pjsua2.NewCallOpParam(false)
	callOpParam.SetStatusCode(call.GetInfo().GetLastStatusCode())
	call.Hangup(callOpParam)
}

func (call *Call) Close() {
	fmt.Println("CLOSE CALL")
	CheckThread()

	call.stopListenAudioSteam()

	if call.audioRecorder != nil && call.player != nil {
		call.player.StopTransmit(call.audioRecorder)
	}

	call.player = nil
	call.audioRecorder = nil
	call.audMed = nil

	EmitEvent("OnCallClosed", call.CALL_ID, "success")

}

func (call *Call) stopListenAudioSteam() {
	if call.ticker != nil {
		call.ticker.Stop()
	}
}

func (call *Call) startCallingTimer() {
	call.callingTimer = time.NewTimer(time.Second * time.Duration(call.CALLING_TIMEOUT))
	fmt.Printf("CALLING TIMER STARTED! CALLING TIMEOUT=%vsec\n", call.CALLING_TIMEOUT)

	go func() {
		<-call.callingTimer.C
		fmt.Println("CALLING TIMER EXPIRED!")
		call.interrupt()
	}()

}

func (call *Call) OnStreamPreCreate(arg2 pjsua2.OnStreamPreCreateParam) {
	fmt.Printf("OnStreamPreCreate\n")
}

func (call *Call) OnInstantMessage(prm pjsua2.OnInstantMessageParam) {
	fmt.Printf("OnInstantMessage, From: %s, To: %s, Message: %s\n",
		prm.GetFromUri(), prm.GetToUri(), prm.GetMsgBody())
}

func (call *Call) OnInstantMessageStatus(prm pjsua2.OnInstantMessageStatusParam) {
	fmt.Printf("OnInstantMessageStatus\n")
}

func (call *Call) OnStreamCreated(prm pjsua2.OnStreamCreatedParam) {
	fmt.Printf("OnStreamCreated\n")
}

func (call *Call) OnStreamDestroyed(prm pjsua2.OnStreamDestroyedParam) {
	fmt.Printf("OnStreamDestroyed\n")
}

func (call *Call) OnDtmfDigit(prm pjsua2.OnDtmfDigitParam) {
	fmt.Printf("OnDtmfDigit\n")
}

func (call *Call) OnCallTransferRequest(prm pjsua2.OnCallTransferRequestParam) {
	fmt.Printf("OnCallTransferRequest\n")
}

func (call *Call) OnCallTransferStatus(prm pjsua2.OnCallTransferStatusParam) {
	fmt.Printf("OnCallTransferStatus\n")
}

func (call *Call) OnCallReplaceRequest(prm pjsua2.OnCallReplaceRequestParam) {
	fmt.Printf("OnCallReplaceRequest\n")
}

func (call *Call) OnCallReplaced(prm pjsua2.OnCallReplacedParam) {
	fmt.Printf("OnCallReplaced\n")
}

func (call *Call) OnCallMediaTransportState(prm pjsua2.OnCallMediaTransportStateParam) {
	fmt.Printf("OnCallMediaTransportState\n")
}

func (call *Call) OnCreateMediaTransport(prm pjsua2.OnCreateMediaTransportParam) {
	fmt.Printf("OnCreateMediaTransport\n")
}

func (call *Call) OnCreateMediaTransportSrtp(prm pjsua2.OnCreateMediaTransportSrtpParam) {
	fmt.Printf("OnCreateMediaTransportSrtp\n")
}

func (call *Call) OnCallMediaEvent(prm pjsua2.OnCallMediaEventParam) {
	fmt.Printf("OnCallMediaEvent\n")
}

func (call *Call) OnCallTsxState(prm pjsua2.OnCallTsxStateParam) {
	fmt.Printf("OnCallTsxState, lastStatusCode=%v\n", call.GetInfo().GetLastStatusCode())
}

func (call *Call) OnCallRxReinvite(prm pjsua2.OnCallRxReinviteParam) {
	fmt.Printf("OnCallRxReinvite\n")
}

func (call *Call) OnCallRxOffer(prm pjsua2.OnCallRxOfferParam) {
	fmt.Printf("OnCallRxOffer\n")
}

func (call *Call) OnCallTxOffer(prm pjsua2.OnCallTxOfferParam) {
	fmt.Printf("OnCallTxOffer\n")
}

func (call *Call) OnTypingIndication(prm pjsua2.OnTypingIndicationParam) {
	fmt.Printf("OnTypingIndication\n")
}

func (call *Call) OnCallRedirected(prm pjsua2.OnCallRedirectedParam) {
	fmt.Printf("OnCallRedirected\n")
}
