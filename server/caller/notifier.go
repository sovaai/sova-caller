package caller

import (
	"fmt"
	"sova-caller-backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Notification struct {
	CLIENT_ID string `json:"id"`
}

var Notifier chan *Notification

func init() {
	Notifier = make(chan *Notification)
}

func Stop(ctx *fiber.Ctx) error {
	defer utils.TimeTrack(time.Now(), "NOTIFY")

	data := new(Notification)

	if err := ctx.BodyParser(data); err != nil {
		fmt.Println("error = ", err)
		return ctx.SendStatus(400)
	}

	go Notify(data)

	return ctx.SendStatus(200)
}

func Notify(data *Notification) {
	Notifier <- data
}
