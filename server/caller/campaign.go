package caller

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sova-caller-backend/api"
	"sova-caller-backend/constants"
	"sova-caller-backend/sip"
	"sova-caller-backend/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/looplab/fsm"
)

/* TYPES */

type CampaignData struct {
	UUID      string             `json:"uuid"`
	CAPI      string             `json:"capi"`
	CLIENT_ID string             `json:"id,omitempty"`
	CONTACTS  []string           `json:"contacts"`
	SIP       sip.SIPCredentials `json:"sip"`
	ASR       api.ServiceConfig  `json:"asr"`
	TTS       api.ServiceConfig  `json:"tts"`
}

type CampaignState struct {
	SESSION_ID  string
	CAMPAIGN_ID string
	CLIENT_ID   string
	UUID        string
	CAPI        string
	CONTACTS    []string
	SIPUSER     *SipUser
	ACCOUNT     *sip.Account
	FSM         *fsm.FSM
	CURRENT     int
	SIP         sip.SIPCredentials
	ASR         api.ServiceConfig
	TTS         api.ServiceConfig
}

type SessionData struct {
	CONTACT string
	ACCOUNT *sip.Account
	UUID    string
}

type SessionEntity struct {
	DATA struct {
		ID json.Number `json:"id"`
	} `json:"data"`
}

/* TYPES */

func CreateCampaignState(data *CampaignData) *CampaignState {

	state := &CampaignState{
		SESSION_ID:  "0",
		UUID:        data.UUID,
		CAPI:        data.CAPI,
		CONTACTS:    data.CONTACTS,
		CAMPAIGN_ID: uuid.New().String(),
		CLIENT_ID:   data.CLIENT_ID,
		SIPUSER:     nil,
		ACCOUNT:     nil,
		CURRENT:     0,
		SIP:         data.SIP,
		ASR:         data.ASR,
		TTS:         data.TTS,
	}

	state.FSM = fsm.NewFSM(
		"init",
		fsm.Events{
			{Name: "run", Src: []string{"init", "finished"}, Dst: "running"},
			{Name: "stop", Src: []string{"running"}, Dst: "aborted"},
			{Name: "end", Src: []string{"running"}, Dst: "finished"},
		},
		fsm.Callbacks{},
	)

	return state
}

func PrepareSessionData(campaign *CampaignState) SessionData {
	return SessionData{CONTACT: campaign.CONTACTS[campaign.CURRENT], UUID: campaign.UUID, ACCOUNT: campaign.ACCOUNT}
}

func CreateCampaignDir(CAMPAIGN_ID string) {
	dirPath := filepath.Join(constants.AUDIO_DIR, CAMPAIGN_ID, constants.SESSIONS_DIR)
	if _, err := os.Stat(dirPath); err != nil {
		merr := os.MkdirAll(dirPath, os.ModePerm)
		if merr != nil {
			fmt.Println(merr)
		}
	}
}

func CampaignWorkflow(data *CampaignData) {
	defer func() {

		if err := recover(); err != nil {
			fmt.Println("CAMPAIGN ERROR", err)
		} else {
			fmt.Println("CAMPAIGN ROUTINE WAS FINISHED!")
		}

	}()

	campaign := CreateCampaignState(data)

	status := make(chan string)

	campaign.SIPUSER = CreateCaller(status, campaign.UUID)
	campaign.ACCOUNT = sip.RegisterAccount(campaign.SIP, campaign.UUID)

	CreateCampaignDir(campaign.CAMPAIGN_ID)

	campaign.SIPUSER.StartCaller()

	for {
		select {
		case statusNotification := <-Notifier:

			if statusNotification != nil && campaign.CLIENT_ID == statusNotification.CLIENT_ID {
				campaign.FSM.Event("stop")
			}

		case current := <-status:
			if current == "OnSipReady" {

				campaign.FSM.Event("run")
				ClientReporter.UpdateStatus(Report{
					CLIENT_ID:       campaign.CLIENT_ID,
					CAMPAIGN_STATUS: campaign.FSM.Current(),
				})
			}

			if current == "OnRegState" {

				Session(PrepareSessionData(campaign), campaign)
			}

			if current == "OnCallClosed" {
				fmt.Println("CAMPAIGN ON CALL CLOSED!")

				if campaign.FSM.Current() == "aborted" {
					sip.EmitEvent("OnUnregister")
					continue
				}

				campaign.CURRENT += 1

				if campaign.CURRENT < len(campaign.CONTACTS) {
					Session(PrepareSessionData(campaign), campaign)
					continue
				}

				fmt.Println("SHOULD FINISH")
				campaign.FSM.Event("end")
				sip.EmitEvent("OnUnregister")
			}

			if current == "OnUnregister" {
				fmt.Println("FINISH PJSUA CAMPAIGN")
				sip.CheckThread()

				campaignStatus := utils.If(!campaign.ACCOUNT.GetInfo().GetRegIsActive(), "failed", campaign.FSM.Current())

				ClientReporter.UpdateStatus(Report{
					CLIENT_ID:       campaign.CLIENT_ID,
					CAMPAIGN_STATUS: campaignStatus,
				})

				campaign.ACCOUNT.Shutdown()

				sip.Handlers = make(map[sip.IUserService]bool)
				sip.Calls = make(map[string]*sip.Call)

				close(status)

				return
			}
		}
	}
}

func Session(payload SessionData, campaign *CampaignState) {

	defer func() {
		if err := recover(); err != nil {
			sip.EmitEvent("OnCallClosed", fmt.Sprintf("%v", payload.CONTACT), "error")
			fmt.Println("Panic occurred in SESSION FUNCTION:", err)
		}
	}()

	fmt.Println("INIT NEXT SESSION")

	sessionId := strconv.FormatInt(time.Now().UnixMilli(), 10)

	data := &sip.MakeCallData{
		FROM_UID:    campaign.SIP.SIP_ID,
		TO_UID:      payload.CONTACT,
		UUID:        payload.UUID,
		ACCOUNT:     payload.ACCOUNT,
		SESSION_ID:  sessionId,
		CAPI:        campaign.CAPI,
		CAMPAIGN_ID: campaign.CAMPAIGN_ID,
		REGISTRAR:   campaign.SIP.REGISTRAR,
		ASR:         campaign.ASR,
		TTS:         campaign.TTS,
	}

	sip.MakeCall(data)
}

func Workflow(c *fiber.Ctx) error {
	defer utils.TimeTrack(time.Now(), "CAMPAIGN")

	data := &CampaignData{}
	if err := c.BodyParser(data); err != nil {
		fmt.Println("error = ", err)
		return c.SendStatus(400)
	}

	if data.CLIENT_ID == "" {
		data.CLIENT_ID = uuid.New().String()
	}

	go CampaignWorkflow(data)

	return c.SendString(data.CLIENT_ID)
}
