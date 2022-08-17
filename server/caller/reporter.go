package caller

import (
	"fmt"
	"sova-caller-backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Report struct {
	CLIENT_ID       string `json:"client_id"`
	CAMPAIGN_STATUS string `json:"campaign_status,omitempty"`
}

type Reporter struct {
	Clients map[string]string
}

var ClientReporter *Reporter

func init() {
	ClientReporter = NewReporter()
}

func NewReporter() *Reporter {
	reporter := &Reporter{
		Clients: make(map[string]string),
	}

	return reporter
}

func (reporter *Reporter) UpdateStatus(data Report) {
	reporter.Clients[data.CLIENT_ID] = data.CAMPAIGN_STATUS

	fmt.Printf("CLIENT_ID: %v, CAMPAIGN_STATUS: %v\n",
		data.CLIENT_ID, data.CAMPAIGN_STATUS)
}

func GetStatus(ctx *fiber.Ctx) error {
	defer utils.TimeTrack(time.Now(), "GET CAMPAIGN STATUS")

	client_id := ctx.Params("client_id")

	fmt.Printf("GET CAMPAIGN STATUS: CLIENT_ID: %v\n", client_id)

	if ClientReporter.Clients[client_id] == "" {
		return ctx.SendStatus(404)
	}

	return ctx.SendString(ClientReporter.Clients[client_id])
}
