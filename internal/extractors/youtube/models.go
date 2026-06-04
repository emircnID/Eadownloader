package youtube

type Info struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Uploader    string   `json:"uploader"`
	Duration    float64  `json:"duration"`
	Thumbnail   string   `json:"thumbnail"`
	WebpageURL  string   `json:"webpage_url"`
	Formats     []Format `json:"formats"`
	RequestedID string   `json:"-"`
}

type Format struct {
	FormatID       string  `json:"format_id"`
	URL            string  `json:"url"`
	Ext            string  `json:"ext"`
	Protocol       string  `json:"protocol"`
	VideoCodec     string  `json:"vcodec"`
	AudioCodec     string  `json:"acodec"`
	Width          int32   `json:"width"`
	Height         int32   `json:"height"`
	TBR            float64 `json:"tbr"`
	Filesize       int64   `json:"filesize"`
	FilesizeApprox int64   `json:"filesize_approx"`
}
