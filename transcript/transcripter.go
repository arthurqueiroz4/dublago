package transcript

type Transcripter interface {
	Transcript(fPath string) (string, error)
}
