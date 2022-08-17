package main

import (
	"log"
	"sova-caller-backend/caller"
	"sova-caller-backend/constants"
	"sova-caller-backend/sip"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	if sip.Endpoint == nil {
		sip.InitLib()
	}
}

func main() {
	PORT := constants.EnvVariable("PORT")

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(helmet.New())
	app.Use(etag.New(etag.Config{
		Weak: true,
	}))
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET,POST,OPTIONS",
		ExposeHeaders: "Content-Type,Authorization,Accept",
	}))
	app.Use(requestid.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Get("/status/:client_id", caller.GetStatus)

	app.Use("/campaign/start", caller.Workflow)
	app.Use("/campaign/stop", caller.Stop)

	app.Static("/", "public")
	app.Static("/audio", "audio", fiber.Static{
		ByteRange: true,
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	log.Fatal(app.Listen(":" + PORT))
}
