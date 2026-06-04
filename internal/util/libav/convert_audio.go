package libav

import (
	"fmt"
	"os"
	"os/exec"
)

func ConvertAudioToMP3(inputPath string, outputPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", inputPath,
		"-vn",
		"-codec:a", "libmp3lame",
		"-q:a", "2",
		outputPath,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		os.Remove(outputPath)
		return fmt.Errorf("failed to convert audio to mp3: %w: %s", err, string(output))
	}
	return nil
}
