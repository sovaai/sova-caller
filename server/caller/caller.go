package caller

import (
	"fmt"
	"log"
	"sova-caller-backend/sip"
)

type SipUser struct {
	callId    string
	state     string
	stateChan chan string
	uuid      string
	reason    string
}

func CreateCaller(state chan string, uuid string) *SipUser {
	return &SipUser{stateChan: state, uuid: uuid, reason: ""}
}

func (su *SipUser) StartCaller() {

	sip.RegisterEventHandler(su)
}

func (su *SipUser) OnSipReady() {
	su.state = "OnSipReady"
	su.stateChan <- su.state
}

func (su *SipUser) OnRegState(uid string, isActive bool, code int) {
	fmt.Printf("OnRegState, userId = %s, isActive = %v, code = %v\n", uid, isActive, code)

	if isActive {
		su.state = "OnRegState"
	}

	su.stateChan <- su.state
}

func (su *SipUser) OnCallClosed(callId string, reason string) {
	log.Println("Caller - Close call")
	su.state = "OnCallClosed"
	su.stateChan <- su.state
	su.reason = reason
}

func (su *SipUser) GetReason() string {
	return su.reason
}

func (su *SipUser) OnUnregister() {
	su.state = "OnUnregister"
	su.stateChan <- su.state
}
