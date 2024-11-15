package tts

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoquittsWithDefaultModel(t *testing.T) {
	var s Speecher = NewCoquiTTS()
	outPath, err := s.Speech("Mensagem para ser falada pelo TTS")
	assert.Nil(t, err, "err should be nil in speech return")
	assert.NotNil(t, outPath, "outPath should be not nil in Speech return")

	fileContent, err := os.ReadFile(outPath)
	assert.Nil(t, err, "err should be nil in speech return")
	assert.NotNil(t, fileContent, "fileContent should be not nil in Read File")
	assert.NotEmpty(t, fileContent, "fileContent should be not not empty in Read File")

	os.Remove(outPath)
}

func TestCoquittsWithYourTTSModel(t *testing.T) {
	var s Speecher = NewCoquiTTS(
		WithModelName("tts_models/multilingual/multi-dataset/your_tts"),
		WithLanguageIdx("pt-br"),
	)
	outPath, err := s.Speech("Mensagem para ser falada pelo TTS")
	assert.Nil(t, err, "err should be nil in speech return")
	assert.NotNil(t, outPath, "outPath should be not nil in Speech return")

	fileContent, err := os.ReadFile(outPath)
	assert.Nil(t, err, "err should be nil in speech return")
	assert.NotNil(t, fileContent, "fileContent should be not nil in Read File")
	assert.NotEmpty(t, fileContent, "fileContent should be not not empty in Read File")

	os.Remove(outPath)
}
