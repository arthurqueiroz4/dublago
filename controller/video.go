package controller

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
)

func SendVideo(c fiber.Ctx) error {
	id := c.Locals("user.sso_id").(string)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if fileHeader == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "video not found"})
	}
	videosPath, _ := filepath.Abs("videos/")
	userVideosPath := videosPath + "/" + id
	err = os.MkdirAll(userVideosPath, 0755)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c.SaveFile(fileHeader, userVideosPath+".mp4")

	// path, _ := filepath.Abs("pipe/")
	// singleton.Pipeline.Run(username)

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
