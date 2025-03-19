package tts

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CoquiTTS struct {
	modelName      string
	speakerWavPath string
	languageIdx    string
}

type Option func(*CoquiTTS)

func NewCoquiTTS(opts ...Option) *CoquiTTS {
	cq := new(CoquiTTS)

	cq.modelName = "tts_models/multilingual/multi-dataset/xtts_v2"
	cq.speakerWavPath, _ = filepath.Abs("pipe/primegean-speaker.wav")
	cq.languageIdx = "pt"

	for _, o := range opts {
		o(cq)
	}

	return cq
}

func WithModelName(modelName string) Option {
	return func(ct *CoquiTTS) {
		ct.modelName = modelName
	}
}

func WithSpeakerWavPath(speakerWavRelativePath string) Option {
	return func(ct *CoquiTTS) {
		ct.speakerWavPath, _ = filepath.Abs(speakerWavRelativePath)
	}
}

func WithLanguageIdx(languageIdx string) Option {
	return func(ct *CoquiTTS) {
		ct.languageIdx = languageIdx
	}
}

func (c *CoquiTTS) Speech(message, speakerWavPath string) (string, error) {
	slog.Info("[TTS] Generating speech now")

	t := time.Now()
	outFile := fmt.Sprintf("%d-%d-%d-%s-tts.wav", t.Day(), t.Month(), t.Year(), uuid.New())
	outPath, err := filepath.Abs("pipe/" + outFile)
	if err != nil {
		slog.Error(`[TTS] Error in getting absolute outPath path: ` + err.Error())
		return "", err
	}

	text := strings.ReplaceAll(message, ".", `\n`)
	slog.Info("[TTS] Model Name: " + c.modelName)
	slog.Info("[TTS] Text: " + text)
	slog.Info("[TTS] Speaker Wav Path: " + speakerWavPath)
	slog.Info("[TTS] Out Path: " + outPath)
	slog.Info("[TTS] Language Index: " + c.languageIdx)

	cmd := exec.Command("tts",
		"--model_name", c.modelName,
		"--text", text,
		"--speaker_wav", speakerWavPath,
		"--out_path", outPath,
		"--language_idx", c.languageIdx)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		slog.Error(`[TTS] Error in running 'cmd.Run': ` + err.Error() +
			"\n\t\t\tCLI Error Output: " + stderr.String())
		return "", err
	}

	slog.Info(`[TTS] Speech generated with success
	Output path: ` + outPath)

	return outPath, nil
}
