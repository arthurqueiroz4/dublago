package main

import (
	"log"
	"tradutor-dos-crias/auth"
	"tradutor-dos-crias/caption"
	"tradutor-dos-crias/media"
	"tradutor-dos-crias/pipeline"
	"tradutor-dos-crias/transcript"
	"tradutor-dos-crias/translator"
	"tradutor-dos-crias/tts"

	"github.com/gofiber/fiber/v3"
)

var Pipeline *pipeline.Pipeline

func main() {
	mediaHandler := &media.FfmpegWrapper{}
	transcripter := transcript.NewWhisper()
	translator := &translator.GoogleTranslator{}
	tts := tts.NewCoquiTTS()
	subtitler := caption.NewStablets()

	Pipeline = pipeline.NewPipeline(transcripter, mediaHandler, translator, tts, subtitler)

	app := fiber.New()

	app.Get("/api/auth/callback", auth.Callback)

	log.Fatal(app.Listen(":4000"))
}
