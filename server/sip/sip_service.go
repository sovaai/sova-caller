package sip

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"sova-caller-backend/api"

	pjsua2 "sova-caller-backend/pjsua2"
)

/* TYPES */

type MakeCallData struct {
	FROM_UID    string
	TO_UID      string
	UUID        string
	ACCOUNT     *Account
	SESSION_ID  string
	CAPI        string
	CAMPAIGN_ID string
	REGISTRAR   string
	ASR         api.ServiceConfig
	TTS         api.ServiceConfig
}

type SIPCredentials struct {
	REGISTRAR    string `json:"sip_registrar"`
	SIP_ID       string `json:"sip_id"`
	SIP_PASSWORD string `json:"sip_password"`
}

/* TYPES */

var (
	Endpoint pjsua2.Endpoint
	Handlers = make(map[IUserService]bool)

	mutex sync.Mutex
	Calls map[string]*Call
)

func InitLib() {

	Endpoint = pjsua2.NewEndpoint()

	Calls = make(map[string]*Call)

	// Create endpoint
	Endpoint.LibCreate()

	// Initialize endpoint
	epConfig := pjsua2.NewEpConfig()
	epConfig.GetUaConfig().SetUserAgent(config.AgentName)
	epConfig.GetUaConfig().SetMaxCalls(config.MaxCall)

	epConfig.GetLogConfig().SetLevel(config.PjLogLevel)
	if config.PjLogLevel > 0 {
		epConfig.GetLogConfig().SetWriter(pjsua2.NewDirectorLogWriter(new(LogWriter)))
	}

	epConfig.GetMedConfig().SetNoVad(false)

	Endpoint.LibInit(epConfig)
	Endpoint.AudDevManager().SetNullDev()
	Endpoint.AudDevManager().SetVad(true, false)

	transportConfig := pjsua2.NewTransportConfig()
	transportConfig.SetPort(config.LocalPort)

	var transport pjsua2.Pjsip_transport_type_e
	if strings.EqualFold(config.Transport, "UDP") {
		transport = pjsua2.PJSIP_TRANSPORT_UDP
	} else if strings.EqualFold(config.Transport, "TLS") {
		transport = pjsua2.PJSIP_TRANSPORT_TLS
	} else if strings.EqualFold(config.Transport, "TCP") {
		transport = pjsua2.PJSIP_TRANSPORT_TCP
	} else {
		fmt.Printf("unknown config.Transport = %s\n", config.Transport)
	}

	Endpoint.TransportCreate(transport, transportConfig)
	Endpoint.LibStart()

	fmt.Printf("[ sip.Service ] Available codecs:\n")
	for i := 0; i < int(Endpoint.CodecEnum2().Size()); i++ {
		c := Endpoint.CodecEnum2().Get(i)
		fmt.Printf("\t - %s (priority: %d)\n", c.GetCodecId(), c.GetPriority())
		if !strings.HasPrefix(c.GetCodecId(), "PCM") {
			Endpoint.CodecSetPriority(c.GetCodecId(), 0)
		}
	}

	fmt.Printf("[ sip.Service ] PJSUA2 STARTED ***\n")
}

func getRemoteURI(remoteUser string, registrar string) string {

	CheckThread()

	remoteUri := strings.Builder{}

	remoteUri.WriteString("sip:")
	remoteUri.WriteString(remoteUser)
	remoteUri.WriteString("@")
	remoteUri.WriteString(registrar)
	remoteUri.WriteString(":")
	remoteUri.WriteString(fmt.Sprintf("%d", config.ProxyPort))
	if !strings.EqualFold(config.Transport, "UDP") {
		if strings.EqualFold(config.Transport, "TCP") {
			remoteUri.WriteString(";transport=tcp")
		} else {
			remoteUri.WriteString(";transport=tls")
		}
	}
	return remoteUri.String()
}

func EmitEvent(event string, params ...interface{}) {
	// emit event at new thread
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("panic occurred in sip_service EmitEvent method:", err, event)
			}
		}()

		for h := range Handlers {
			switch event {
			case "OnSipReady":
				h.OnSipReady()
			case "OnRegState":
				h.OnRegState(params[0].(string), params[1].(bool), params[2].(int))
			case "OnCallClosed":
				h.OnCallClosed(params[0].(string), params[1].(string))
			case "OnUnregister":
				h.OnUnregister()
			default:
				fmt.Printf("EmitEvent, unknown event = %s\n", event)
			}
		}
	}()
}

func onRegState(uri string, active bool, code pjsua2.Pjsip_status_code) {
	CheckThread()
	EmitEvent("OnRegState", uri, active, int(code))
}

func OnCallClosed(callId string, reason string) {
	delete(Calls, callId)
	EmitEvent("OnCallClosed", callId, reason)
}

func CheckThread() {
	mutex.Lock()
	defer mutex.Unlock()

	if !Endpoint.LibIsThreadRegistered() {
		Endpoint.LibRegisterThread(config.AgentName)
	}
}

func RegisterEventHandler(sipUser IUserService) {
	Handlers[sipUser] = true
	CheckThread()
	EmitEvent("OnSipReady")
}

func RegisterAccount(data SIPCredentials, uuid string) *Account {
	CheckThread()

	idUri := fmt.Sprintf("sip:%s@%s", data.SIP_ID, data.REGISTRAR)
	registrarUri := fmt.Sprintf("sip:%s", data.REGISTRAR)

	accountConfig := pjsua2.NewAccountConfig()

	cred := pjsua2.NewAuthCredInfo("digest", "*", data.SIP_ID, 0, data.SIP_PASSWORD)

	aSipCfg := pjsua2.NewAccountSipConfig()

	accountConfig.SetSipConfig(aSipCfg)
	accountConfig.SetIdUri(idUri)
	accountConfig.GetRegConfig().SetRegistrarUri(registrarUri)
	accountConfig.GetSipConfig().GetAuthCreds().Add(cred)

	account := NewAccount(uuid)

	pjAccount := pjsua2.NewDirectorAccount(account)

	pjAccount.Create(accountConfig)

	account.Account = pjAccount
	fmt.Printf("CREATE LOCAL ACCOUNT: Account Created, URI=%s\n", account.GetInfo().GetUri())

	return account
}

func MakeCall(data *MakeCallData) {
	CheckThread()

	call := NewCall(&NewCallData{
		UUID:        data.UUID,
		CAMPAIGN_ID: data.CAMPAIGN_ID,
		SESSION_ID:  data.SESSION_ID,
	})

	call.dosWorkflow = api.NewDosWorkflow(&api.NewDosWorkflowData{
		CAPI:        data.CAPI,
		UUID:        call.uuid,
		SESSION_ID:  data.SESSION_ID,
		CAMPAIGN_ID: data.CAMPAIGN_ID,
		ASR:         data.ASR,
		TTS:         data.TTS,
	})

	call.asrResponseReceived = make(chan string)
	call.needClose = make(chan string)
	call.error = make(chan string)
	call.ticker = time.NewTicker(CHECK_STREAM_INTERVAL_MSEC * time.Millisecond)

	remoteUri := getRemoteURI(data.TO_UID, data.REGISTRAR)

	CheckThread()

	call.Call = pjsua2.NewDirectorCall(call, data.ACCOUNT)

	callOpParam := pjsua2.NewCallOpParam(true)
	callSetting := pjsua2.NewCallSetting()
	callSetting.SetAudioCount(1)
	callOpParam.SetOpt(callSetting)

	fmt.Printf("MakeCall, From=%s, To=%s\n", data.ACCOUNT.GetInfo().GetUri(), remoteUri)

	CheckThread()

	call.MakeCall(remoteUri, callOpParam)

	Calls[call.GetInfo().GetCallIdString()] = call
}
