package input

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

const (
	bestvideo string = "bestvideo"
	bestaudio string = "bestaudio"
)

type YouTube struct{}

func New() *YouTube {
	return &YouTube{}
}

func (y *YouTube) Download(url string) (string, string) {
	slog.Info("[YOUTUBE] Downloading video and audio with 'yt-dlp'")

	path, _ := filepath.Abs("pipe/")
	videoPath := path + "/video.mp4"
	audioPath := path + "/audio.mp4"

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go y.downloadStream(url, bestvideo, videoPath, wg)
	go y.downloadStream(url, bestaudio, audioPath, wg)

	wg.Wait()

	return videoPath, audioPath
}

func (y *YouTube) downloadStream(url, format, outputPath string, wg *sync.WaitGroup) {
	slog.Info("[YOUTUBE] Downloading " + format + " by " + url)
	defer func() {
		slog.Info("[YOUTUBE] Download finished: " + format + " by " + url)
		wg.Done()
	}()

	cmd := exec.Command("yt-dlp", []string{
		"-f", format,
		"-o", outputPath,
		url,
	}...)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		os.Stdout.WriteString("[YOUTUBE] yt-dlp stderr:\n\t" + stderr.String())
		panic(err)
	}
}
