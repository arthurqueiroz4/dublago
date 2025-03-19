package pipeline

import (
	"os"
	"path/filepath"
	"tradutor-dos-crias/caption"
	"tradutor-dos-crias/input"
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
	subtitler    caption.Subtitler
	downloader   input.Downloader
}

func NewPipeline(transcripter transcript.Transcripter,
	mediaHandler media.MediaHandler,
	translator translator.Translator,
	speaker tts.Speecher,
	subtitler caption.Subtitler,
	downloader input.Downloader,
) *Pipeline {
	i := &Pipeline{
		transcripter,
		mediaHandler,
		translator,
		speaker,
		subtitler,
		downloader,
	}

	return i
}

func (i *Pipeline) RunWithYoutube(url, outputVideoPath string) error {
	outputFolderDefault, err := filepath.Abs("pipe/")
	if err != nil {
		return err
	}
	filesToRemove := make([]string, 0)
	defer removeTmpFiles(&filesToRemove)

	inputVideoPath, inputAudioPath := i.downloader.Download(url)
	filesToRemove = append(filesToRemove, inputVideoPath, inputAudioPath)

	wavAudioPath := outputFolderDefault + "/" + uuid.NewString() + ".wav"
	filesToRemove = append(filesToRemove, wavAudioPath)

	err = i.mediaHandler.ConvertToWav(inputAudioPath, wavAudioPath)
	if err != nil {
		return err
	}

	text, err := i.transcripter.Transcript(inputAudioPath)
	filesToRemove = append(filesToRemove, "audio.txt")
	if err != nil {
		return err
	}

	translatedText, err := i.translator.Translate(text)
	if err != nil {
		return err
	}

	dubbedAudio, err := i.speaker.Speech(translatedText, wavAudioPath)
	if err != nil {
		return err
	}

	tmpDubbedFileName := outputFolderDefault + "/" + uuid.NewString() + ".mp4"
	filesToRemove = append(filesToRemove, tmpDubbedFileName, dubbedAudio)

	err = i.mediaHandler.Merge(inputVideoPath, dubbedAudio, tmpDubbedFileName)
	if err != nil {
		return err
	}

	subtitlesPath, err := i.subtitler.GenerateSubtitle(tmpDubbedFileName)
	if err != nil {
		return err
	}
	filesToRemove = append(filesToRemove, subtitlesPath)
	_, err = i.mediaHandler.MergeSubtitle(tmpDubbedFileName, subtitlesPath, outputVideoPath)
	if err != nil {
		return err
	}

	return nil
}

func (i *Pipeline) RunWithLocalVideo(videoPath, outputVideoPath string) error {
	filesToRemove := make([]string, 0)
	filesToRemove = append(filesToRemove, videoPath)
	defer removeTmpFiles(&filesToRemove)

	outputFolderDefault, err := filepath.Abs("pipe/")
	if err != nil {
		return err
	}

	unmergedVideoPath := outputFolderDefault + "/" + uuid.NewString() + ".mp4"
	unmergedAudioPath := outputFolderDefault + "/" + uuid.NewString() + ".mp4"
	filesToRemove = append(filesToRemove, unmergedVideoPath, unmergedAudioPath)

	err = i.mediaHandler.Unmerge(videoPath, unmergedVideoPath, unmergedAudioPath)
	if err != nil {
		return err
	}

	text, err := i.transcripter.Transcript(unmergedAudioPath)
	filesToRemove = append(filesToRemove, "audio.txt")
	if err != nil {
		return err
	}

	translatedText, err := i.translator.Translate(text)
	if err != nil {
		return err
	}

	unmergedAudioWavPath := outputFolderDefault + "/" + uuid.NewString() + ".mp4"
	filesToRemove = append(filesToRemove, unmergedAudioPath)

	err = i.mediaHandler.ConvertToWav(unmergedAudioPath, unmergedAudioWavPath)
	if err != nil {
		return err
	}

	dubbedAudio, err := i.speaker.Speech(translatedText, unmergedAudioWavPath)
	if err != nil {
		return err
	}
	filesToRemove = append(filesToRemove, dubbedAudio)

	err = i.mediaHandler.Merge(unmergedVideoPath, dubbedAudio, outputVideoPath)
	if err != nil {
		return err
	}

	// subtitlesPath, err := i.subtitler.GenerateSubtitle(outputVideoPath)
	// if err != nil {
	// 	return err
	// }
	// filesToRemove = append(filesToRemove, subtitlesPath)
	// _, err = i.mediaHandler.MergeSubtitle(outputVideoPath, subtitlesPath, outputVideoPath)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func removeTmpFiles(filesPath *[]string) {
	for _, file := range *filesPath {
		os.Remove(file)
	}
}
