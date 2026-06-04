package youtube

type PlayerResponse struct {
	PlayabilityStatus struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	} `json:"playabilityStatus"`
	VideoDetails struct {
		Title         string `json:"title"`
		LengthSeconds string `json:"lengthSeconds"`
		Thumbnail     struct {
			Thumbnails []Thumbnail `json:"thumbnails"`
		} `json:"thumbnail"`
	} `json:"videoDetails"`
	StreamingData struct {
		Formats         []Format `json:"formats"`
		AdaptiveFormats []Format `json:"adaptiveFormats"`
	} `json:"streamingData"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
}

type Format struct {
	Itag              int    `json:"itag"`
	URL               string `json:"url"`
	SignatureCipher   string `json:"signatureCipher"`
	Cipher            string `json:"cipher"`
	MimeType          string `json:"mimeType"`
	Bitrate           int64  `json:"bitrate"`
	Width             int32  `json:"width"`
	Height            int32  `json:"height"`
	ContentLength     string `json:"contentLength"`
	ApproxDurationMS  string `json:"approxDurationMs"`
	QualityLabel      string `json:"qualityLabel"`
	AudioQuality      string `json:"audioQuality"`
	AudioSampleRate   string `json:"audioSampleRate"`
	AudioChannels     int32  `json:"audioChannels"`
}

type playerClient struct {
	Name           string
	Version        string
	AndroidSDK     int
	UserAgent      string
}
