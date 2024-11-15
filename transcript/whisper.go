package transcript

import (
	"bytes"
	"errors"
	"log/slog"
	"os/exec"
	"regexp"
	"strings"
)

var availableModels = [...]string{"turbo", "tiny", "base", "small", "medium", "large"}

const (
	timeStampRegex = `\[.*\]`
	errRegex       = `Skipping.*|Error.*`
)

type Whisper struct {
	model    string
	language string
}

type Option func(*Whisper)

func NewWhisper(options ...Option) *Whisper {
	w := new(Whisper)
	w.model = "turbo"
	w.language = "English"

	for _, option := range options {
		option(w)
	}

	return w
}

func WithModel(model string) Option {
	return func(w *Whisper) {
		for _, availableModel := range availableModels {
			if model == availableModel {
				w.model = model
				return
			}
		}
		w.model = "turbo"
	}
}

func WithLanguage(language string) Option {
	return func(w *Whisper) {
		w.language = language
	}
}

func (w *Whisper) Transcript(fPath string) (string, error) {
	slog.Info("[Transcript] Start transcripting audio")

	args := []string{fPath, "--language", w.language, "--model", w.model}
	cmd := exec.Command("whisper", args...)

	var stdErr bytes.Buffer
	var stdOut bytes.Buffer
	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	cmd.Run()
	return handleWhisperOutput(&stdErr, &stdOut)
}

func handleWhisperOutput(stdErr *bytes.Buffer, stdOut *bytes.Buffer) (string, error) {
	errReg := regexp.MustCompile(errRegex)

	errMsgs := errReg.FindAllString(stdErr.String(), -1)

	if len(errMsgs) > 0 {
		errs := strings.Join(errMsgs, ",")
		slog.Error("[Transcript] Error transcripting audio")
		slog.Error("[Transcript] Whisper error message:\n" + stdErr.String())
		return "", errors.New(errs)
	}

	slog.Info("[Transcript] Successfully transcripted text")
	reg := regexp.MustCompile(timeStampRegex)

	text := stdOut.String()
	return reg.ReplaceAllString(text, ""), nil
}
