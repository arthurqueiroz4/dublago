package tts

type Speecher interface {
	Speech(message string) (string, error)
}
