package extractors

import (
	"eadownloader/internal/extractors/facebook"
	"eadownloader/internal/extractors/instagram"
	"eadownloader/internal/extractors/tiktok"
	"eadownloader/internal/extractors/twitter"
	"eadownloader/internal/extractors/youtube"
	"eadownloader/internal/models"
)

var Extractors = []*models.Extractor{
	facebook.WatchShortExtractor,
	facebook.ShareExtractor,
	facebook.Extractor,
	tiktok.Extractor,
	tiktok.VMExtractor,
	twitter.Extractor,
	twitter.ShortExtractor,
	youtube.Extractor,
	instagram.Extractor,
	instagram.StoriesExtractor,
	instagram.ShareURLExtractor,
}
