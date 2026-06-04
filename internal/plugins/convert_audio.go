package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"eadownloader/internal/database"
	"eadownloader/internal/models"
	"eadownloader/internal/util/libav"
)

var ConvertAudioToMP3 = &models.Plugin{
	ID: "convert_audio_to_mp3",
	RunFunc: func(ctx *models.ExtractorContext, _ *models.MediaItem, format *models.DownloadedFormat) error {
		if format.Format.Type != database.MediaTypeAudio {
			return nil
		}

		filePath := format.FilePath
		outputPath := strings.TrimSuffix(
			filePath,
			filepath.Ext(filePath),
		) + "_converted.mp3"
		ctx.FilesTracker.Add(outputPath)

		if err := libav.ConvertAudioToMP3(filePath, outputPath); err != nil {
			return fmt.Errorf("failed to convert audio to mp3: %w", err)
		}

		if err := os.Rename(outputPath, filePath); err != nil {
			return fmt.Errorf("failed to replace audio file: %w", err)
		}

		format.Format.AudioCodec = database.MediaCodecMp3
		return nil
	},
}
