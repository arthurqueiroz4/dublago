package media

type MediaHandler interface {
	Unmerge(inputVideoPath string, outputVideoPath string, outputAudioPath string) error
	Merge(inputVideoPath string, inputAudioPath string, outputVideoPath string) error
}