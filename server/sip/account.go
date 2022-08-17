package sip

import (
	"fmt"
	pjsua2 "sova-caller-backend/pjsua2"
)

type Account struct {
	pjsua2.Account
	uuid string
}

func NewAccount(uuid string) *Account {
	return &Account{uuid: uuid}
}

func (ac *Account) OnRegState(prm pjsua2.OnRegStateParam) {
	CheckThread()

	info := ac.GetInfo()

	var regiState string

	if info.GetRegIsActive() {
		regiState = "REGISTER"
	} else {
		regiState = "UNREGISTER"
		EmitEvent("OnUnregister")
		return
	}

	fmt.Printf("OnRegState, regiState=%v, code=%v\n", regiState, prm.GetCode())

	onRegState(info.GetUri(), info.GetRegIsActive(), prm.GetCode())
}

func (ac *Account) OnIncomingCall(prm pjsua2.OnIncomingCallParam) {

}

func (ac *Account) OnInstantMessage(prm pjsua2.OnInstantMessageParam) {

}

func (ac *Account) OnInstantMessageStatus(prm pjsua2.OnInstantMessageStatusParam) {

}

func (ac *Account) OnRegStarted(prm pjsua2.OnRegStartedParam) {

}

func (ac *Account) Shutdown() {
	CheckThread()
	ac.Account.Shutdown()
}

func (ac *Account) OnIncomingSubscribe(prm pjsua2.OnIncomingSubscribeParam) {
}

func (ac *Account) OnTypingIndication(prm pjsua2.OnTypingIndicationParam) {
}

func (ac *Account) OnMwiInfo(prm pjsua2.OnMwiInfoParam) {
}
