package tts

type Speecher interface {
	Speech(message string, speakerWavPath string) (string, error)
}
