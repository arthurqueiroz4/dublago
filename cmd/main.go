package main

import (
	"log"
	"tradutor-dos-crias/auth"
	"tradutor-dos-crias/config"
	"tradutor-dos-crias/controller"
	"tradutor-dos-crias/database"
	"tradutor-dos-crias/middleware"
	"tradutor-dos-crias/pipeline"
	"tradutor-dos-crias/singleton"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	config.LoadEnvironment()
	err := database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	singleton.Pipeline = pipeline.NewPipeline(singleton.Transcript, singleton.MediaHandler,
		singleton.Translator, singleton.TTS, singleton.Subtitler, singleton.Downloader)

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/auth/callback", auth.Callback)
	app.Get("/api/videos/:filename", controller.GetVideo)

	privateRoute := app.Group("/api", middleware.Authentication)

	privateRoute.Get("/me", controller.Me, middleware.DefaultAuthorization)
	privateRoute.Post("/videos", controller.SendVideo, middleware.DefaultAuthorization)

	log.Fatal(app.Listen(":4000"))
}
