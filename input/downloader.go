package input

type Downloader interface {
	Download(url string) (string, string)
}
