package caption

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Stablets struct {
	model    string
	language string
}

type Option func(*Stablets)

func WithLanguage(language string) Option {
	return func(s *Stablets) {
		s.language = language
	}
}

func WithModel(model string) Option {
	return func(s *Stablets) {
		s.model = model
	}
}

func NewStablets(options ...Option) *Stablets {
	s := new(Stablets)
	s.language = "pt"
	s.model = "turbo"

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *Stablets) GenerateSubtitle(path string) (string, error) {
	slog.Info("[SUBTITLER] Start generate subtitle")

	outputPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".srt"
	args := []string{"--model", s.model, "--language", s.language, "--word_level", "true", "--segment_level", "true", path, "-o", outputPath}

	// Usar stable-ts do virtualenv se existir
	stableTsPath := "./.venv/bin/stable-ts"
	if _, err := os.Stat(stableTsPath); os.IsNotExist(err) {
		stableTsPath = "stable-ts"
	}

	cmd := exec.Command(stableTsPath, args...)

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		slog.Error("[SUBTITLER] Error generating subtitle: " + stdErr.String())
		return "", err
	}

	slog.Info("[SUBTITLER] Succesfully generate subtitle on path: " + outputPath)

	return outputPath, nil
}
