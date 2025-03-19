package controller

import (
	"log/slog"
	"os"
	"path/filepath"
	"tradutor-dos-crias/singleton"

	"github.com/gofiber/fiber/v3"
)

func SendVideo(c fiber.Ctx) error {
	slog.Info("[VIDEO] Sending video")
	id := c.Locals("user.sso_id").(string)

	videosPath, _ := filepath.Abs("videos/")
	userVideosPath := videosPath + "/" + id

	uploadType := c.FormValue("uploadType")

	switch uploadType {
	case "YOUTUBE":
		url := c.FormValue("url")
		err := os.MkdirAll(userVideosPath, 0755)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		err = singleton.Pipeline.RunWithYoutube(url, userVideosPath+".mp4")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	case "FILE":
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		fileInputPath := userVideosPath + "_input" + ".mp4"
		c.SaveFile(fileHeader, fileInputPath)
		err = singleton.Pipeline.RunWithLocalVideo(fileInputPath, userVideosPath+".mp4")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"videoUrl": "http://localhost:4000/api/videos/" + id + ".mp4",
	})
}

func GetVideo(c fiber.Ctx) error {
	filename := c.Params("filename")
	videosPath, _ := filepath.Abs("videos/")
	userVideosPath := videosPath + "/" + filename

	if _, err := os.Stat(userVideosPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Vídeo não encontrado",
		})
	}

	c.Set("Content-Type", "video/mp4")
	c.Set("Accept-Ranges", "bytes")

	return c.SendFile(userVideosPath)
}
