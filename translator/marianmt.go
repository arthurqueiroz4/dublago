package translator

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

type MarianMT struct{}

func (mt MarianMT) Translate(text string) (string, error) {
	slog.Info("[TRANSLATE] Generating translation now")

	scriptPath, err := filepath.Abs("scripts/translator_marianmt.py")
	if err != nil {
		slog.Error(`[TRANSLATE] Error in getting absolute scriptPath path: ` + err.Error())
		return "", err
	}

	outPath, err := filepath.Abs("pipe/tts-marianmt-result.txt")
	if err != nil {
		slog.Error(`[TRANSLATE] Error in getting absolute outPath path: ` + err.Error())
		return "", err
	}

	// Use python do virtualenv se existir, senão use python do sistema
	pythonPath := "./.venv/bin/python"
	if _, err := os.Stat(pythonPath); os.IsNotExist(err) {
		pythonPath = "python3"
	}

	cmd := exec.Command(pythonPath,
		scriptPath,
		"--text", fmt.Sprintf(`"%s"`, text),
		"--output", outPath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		slog.Error(`[TRANSLATE] Error in running 'cmd.Run': ` + err.Error() +
			"\n\t\t\tCLI Error Output: " + stderr.String())
		return "", err
	}

	b, err := os.ReadFile(outPath)
	if err != nil {
		slog.Error(`[TRANSLATE] Error in opening result file: ` + err.Error())
		return "", err
	}

	if len(b) == 0 {
		slog.Error(`[TRANSLATE] Result file is empty`)
		return "", err
	}

	slog.Info(`[TRANSLATE] Translation generated with success`)
	return string(b), nil
}
