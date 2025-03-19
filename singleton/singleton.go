package singleton

import (
	"tradutor-dos-crias/caption"
	"tradutor-dos-crias/input"
	"tradutor-dos-crias/media"
	"tradutor-dos-crias/pipeline"
	"tradutor-dos-crias/transcript"
	"tradutor-dos-crias/translator"
	"tradutor-dos-crias/tts"
	"tradutor-dos-crias/user"
)

var (
	UserService  = &user.UserService{}
	MediaHandler = &media.FfmpegWrapper{}
	Translator   = &translator.MarianMT{}
	Transcript   = transcript.NewWhisper()
	TTS          = tts.NewCoquiTTS()
	Subtitler    = caption.NewStablets()
	Downloader   = &input.YouTube{}
	Pipeline     *pipeline.Pipeline
)
