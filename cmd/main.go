package main

import (
	"log"
	"log/slog"
	"tradutor-dos-crias/auth"
	"tradutor-dos-crias/caption"
	"tradutor-dos-crias/database"
	"tradutor-dos-crias/media"
	"tradutor-dos-crias/middleware"
	"tradutor-dos-crias/pipeline"
	"tradutor-dos-crias/transcript"
	"tradutor-dos-crias/translator"
	"tradutor-dos-crias/tts"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

var Pipeline *pipeline.Pipeline

func main() {
	err := database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	mediaHandler := &media.FfmpegWrapper{}
	transcripter := transcript.NewWhisper()
	translator := &translator.GoogleTranslator{}
	tts := tts.NewCoquiTTS()
	subtitler := caption.NewStablets()

	Pipeline = pipeline.NewPipeline(transcripter, mediaHandler, translator, tts, subtitler)

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/auth/callback", auth.Callback)

	privateRoute := app.Group("/api", middleware.Authentication)

	privateRoute.Get("", func(c fiber.Ctx) error {
		slog.Info("Reached here")
		return c.Status(200).JSON(fiber.Map{"User": "Auth"})
	})

	log.Fatal(app.Listen(":4000"))
}
