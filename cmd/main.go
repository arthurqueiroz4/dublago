package main

import (
	"log"
	"tradutor-dos-crias/auth"
	"tradutor-dos-crias/caption"
	"tradutor-dos-crias/controller"
	"tradutor-dos-crias/database"
	"tradutor-dos-crias/media"
	"tradutor-dos-crias/middleware"
	"tradutor-dos-crias/pipeline"
	"tradutor-dos-crias/singleton"
	"tradutor-dos-crias/transcript"
	"tradutor-dos-crias/translator"
	"tradutor-dos-crias/tts"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

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

	singleton.Pipeline = pipeline.NewPipeline(transcripter, mediaHandler, translator, tts, subtitler)

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/auth/callback", auth.Callback)
	app.Get("/api/videos/:filename", controller.GetVideo)

	privateRoute := app.Group("/api", middleware.Authentication)

	privateRoute.Get("/me", controller.Me, middleware.DefaultAuthorization)
	privateRoute.Post("/videos", controller.SendVideo, middleware.DefaultAuthorization)

	log.Fatal(app.Listen(":4000"))
}
