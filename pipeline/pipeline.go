package pipeline

import (
	"path/filepath"
	"tradutor-dos-crias/media"
	"tradutor-dos-crias/transcript"
	"tradutor-dos-crias/translator"
	"tradutor-dos-crias/tts"

	"github.com/google/uuid"
)

type Pipeline struct {
	transcripter transcript.Transcripter
	mediaHandler media.MediaHandler
	translator   translator.Translator
	speaker      tts.Speecher
}

func NewPipeline(transcripter transcript.Transcripter,
	mediaHandler media.MediaHandler,
	translator translator.Translator,
	speaker tts.Speecher,
) *Pipeline {
	i := &Pipeline{
		transcripter,
		mediaHandler,
		translator,
		speaker,
	}

	return i
}

func (i *Pipeline) Run(inputVideoPath, outputVideoPath string) error {
	outputFolderDefault, err := filepath.Abs("pipe/")
	if err != nil {
		return err
	}

	tmpVideoName := outputFolderDefault + "/" + uuid.New().String() + ".mp4"
	tmpAudioName := outputFolderDefault + "/" + uuid.New().String() + ".mp4"

	err = i.mediaHandler.Unmerge(inputVideoPath, tmpVideoName, tmpAudioName)
	if err != nil {
		return err
	}

	text, err := i.transcripter.Transcript(tmpAudioName)
	if err != nil {
		return err
	}

	translatedText, err := i.translator.Translate(text)
	if err != nil {
		return err
	}

	dubbedAudio, err := i.speaker.Speech(translatedText)
	if err != nil {
		return err
	}

	err = i.mediaHandler.Merge(tmpVideoName, dubbedAudio, outputVideoPath)
	if err != nil {
		return err
	}

	return nil
}
