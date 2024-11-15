package media

import (
	"bytes"
	"log/slog"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type FfmpegWrapper struct{}

const (
	ffmpegOutRegex   = `\[out#[0-9].*`
	ffmpegErrorRegex = "Error.*"
)

func (w *FfmpegWrapper) Merge(inputVideoPath string, inputAudioPath string, outputVideoPath string) error {
	slog.Info("[MEDIA] Start merge video and audio files")

	args := []string{"-i", inputAudioPath, "-i", inputVideoPath, "-c", "copy", outputVideoPath, "-y"}
	cmd := exec.Command("ffmpeg", args...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		slog.Error("[MEDIA] Error on merge files: " + err.Error())

		ffmpegErrorMessage, hasErrorMessage := handleFfmpegLog(stdErr.String(), ffmpegErrorRegex)
		if hasErrorMessage {
			slog.Error("[MEDIA] Ffmpeg error message: \n" + ffmpegErrorMessage)
		}

		return err
	}

	ffmpegOutMessage, hasMessage := handleFfmpegLog(stdErr.String(), ffmpegOutRegex)
	if hasMessage {
		slog.Info("[MEDIA] ffmpeg output message:\n" + ffmpegOutMessage)
	}

	absOutPath, outPathErr := filepath.Abs(outputVideoPath)
	outValue := absOutPath
	if outPathErr != nil {
		outValue = outputVideoPath
	}

	slog.Info("[MEDIA] Successfully merged files into path " + outValue)

	return err
}

func (w *FfmpegWrapper) Unmerge(inputVideoPath string, outputVideoPath string, outputAudioPath string) error {
	inputValue := inputVideoPath

	if inputAbsValue, error := filepath.Abs(inputVideoPath); error == nil {
		inputValue = inputAbsValue
	}

	slog.Info("[MEDIA] Start unmerge audio and video from " + inputValue)
	args := []string{"-i", inputVideoPath, "-an", "-c:v", "copy", outputVideoPath, "-vn", "-c:a", "copy", outputAudioPath, "-y"}
	cmd := exec.Command("ffmpeg", args...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		slog.Error("[MEDIA] Error on unmerge files: " + err.Error())

		ffmpegErrorMessage, hasErrorMessage := handleFfmpegLog(stdErr.String(), ffmpegErrorRegex)
		if hasErrorMessage {
			slog.Error("[MEDIA] Ffmpeg error message:\n" + ffmpegErrorMessage)
		}

		return err
	}

	ffmpegOutMessage, hasMessage := handleFfmpegLog(stdErr.String(), ffmpegOutRegex)
	if hasMessage {
		slog.Info("[MEDIA] ffmpeg output message: \n" + ffmpegOutMessage)
	}

	slog.Info("[MEDIA] Successfully unmerge audio and video")

	return err
}

func handleFfmpegLog(log string, exp string) (string, bool) {
	r, _ := regexp.Compile(exp)

	matches := r.FindAllString(log, -1)

	if len(matches) == 0 {
		return "", false
	}

	return strings.Join(matches, "\n"), true
}
