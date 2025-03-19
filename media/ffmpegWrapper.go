package media

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type FfmpegWrapper struct{}

const (
	ffmpegOutRegex       = `\[out#[0-9].*`
	ffmpegErrorRegex     = "Error.*"
	ffprobeDurationRegex = `\b\d{2}:\d{2}:\d{2}\.\d{2}\b`
)

func (w *FfmpegWrapper) MergeSubtitle(inputVideoPath, inputSubtitlePath, outputVideoPath string) (string, error) {
	slog.Info("[MEDIA] Merge subtittles to video")

	subtitles := fmt.Sprintf("subtitles=%s", inputSubtitlePath)
	args := []string{"-i", inputVideoPath, "-vf", subtitles, outputVideoPath, "-y"}

	cmd := exec.Command("ffmpeg", args...)

	var stdErr bytes.Buffer
	var stdOut bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		slog.Error("[MEDIA] Error on merge files: " + err.Error())
		ffmpegErrorMessage, hasErrorMessage := handleFfmpegLog(stdErr.String(), ffmpegErrorRegex)
		if hasErrorMessage {
			slog.Error("[MEDIA] Ffmpeg error message: \n" + ffmpegErrorMessage)
		}
	}
	slog.Info("[MEDIA] Successfully merge subtitle to video")
	return outputVideoPath, nil
}

func (w *FfmpegWrapper) Merge(inputVideoPath, inputAudioPath, outputVideoPath string) error {
	slog.Info("[MEDIA] Start merge video and audio files")

	speedUpAudioRatio, err := getAudioSpeedRatio(inputVideoPath, inputAudioPath)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	atempo := fmt.Sprintf("atempo=%f", speedUpAudioRatio)

	args := []string{"-i", inputAudioPath, "-i", inputVideoPath, "-filter:a", atempo, outputVideoPath, "-y"}
	cmd := exec.Command("ffmpeg", args...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err = cmd.Run()
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

func (w *FfmpegWrapper) Unmerge(inputVideoPath, outputVideoPath, outputAudioPath string) error {
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

func getAudioSpeedRatio(vPath, aPath string) (float64, error) {
	videoDuration, err := getMediaDuration(vPath)
	if err != nil {
		return 0, err
	}

	audioDuration, err := getMediaDuration(aPath)
	if err != nil {
		return 0, err
	}

	return audioDuration / videoDuration, nil
}

func getMediaDuration(fPath string) (float64, error) {
	args := []string{"-i", fPath}
	cmd := exec.Command("ffprobe", args...)

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		slog.Info(cmd.String())
		return 0, errors.New(stdErr.String())
	}

	outText := stdErr.String()
	r := regexp.MustCompile(ffprobeDurationRegex)
	matches := r.FindAllString(outText, -1)

	if len(matches) == 0 {
		return 0, errors.New("Failed to retrieve duration of file: " + fPath)
	}
	multiplier := []float64{3600, 60, 1}
	duration := strings.Split(matches[0], ":")

	var secondsDuration float64
	for i := 0; i < 3; i++ {
		sTime := duration[i]
		currentDuration, _ := strconv.ParseFloat(sTime, 32)
		secondsDuration += currentDuration * multiplier[i]
	}

	return secondsDuration, nil
}

func (w *FfmpegWrapper) ConvertToWav(inputAudioPath string, wavAudioPath string) error {
	args := []string{"-i", inputAudioPath, "-f", "wav", wavAudioPath}
	cmd := exec.Command("ffmpeg", args...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		slog.Error("[MEDIA] Error on convert to wav: " + err.Error())
		ffmpegErrorMessage, hasErrorMessage := handleFfmpegLog(stdErr.String(), ffmpegErrorRegex)
		if hasErrorMessage {
			slog.Error("[MEDIA] Ffmpeg error message: \n" + ffmpegErrorMessage)
		}

		return err
	}

	slog.Info("[MEDIA] Successfully converted to wav on path " + wavAudioPath)
	return nil
}
