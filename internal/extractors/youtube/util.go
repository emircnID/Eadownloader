package youtube

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	"eadownloader/internal/models"
	"eadownloader/internal/networking"
	"eadownloader/internal/util"
)

var (
	innertubeKeyPattern      = regexp.MustCompile(`"INNERTUBE_API_KEY"\s*:\s*"([^"]+)"`)
	initialPlayerRespPattern = regexp.MustCompile(`ytInitialPlayerResponse\s*=\s*(\{.+?\});`)
)

var playerClients = []playerClient{
	{
		Name:       "ANDROID",
		Version:    "19.09.37",
		AndroidSDK: 30,
		UserAgent: "com.google.android.youtube/19.09.37 (Linux; U; Android 11) gzip",
	},
	{
		Name:      "WEB",
		Version:   "2.20240221.08.00",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	},
}

func GetPlayerResponse(ctx *models.ExtractorContext) (*PlayerResponse, error) {
	watchHTML, apiKey, err := fetchWatchPage(ctx)
	if err != nil {
		return nil, err
	}

	if resp := parseInitialPlayerResponse(watchHTML); hasUsableFormats(resp) {
		return resp, nil
	}

	for _, client := range playerClients {
		resp, err := fetchPlayerAPI(ctx, apiKey, client)
		if err != nil {
			ctx.Warnf("youtube player client %s failed: %v", client.Name, err)
			continue
		}
		if hasUsableFormats(resp) {
			return resp, nil
		}
		if resp.PlayabilityStatus.Status != "" && resp.PlayabilityStatus.Status != "OK" {
			return nil, playabilityError(resp)
		}
	}

	if resp := parseInitialPlayerResponse(watchHTML); resp != nil {
		return nil, playabilityError(resp)
	}

	return nil, fmt.Errorf("no usable youtube formats found")
}

func fetchWatchPage(ctx *models.ExtractorContext) ([]byte, string, error) {
	resp, err := ctx.Fetch(
		http.MethodGet,
		"https://www.youtube.com/watch?v="+url.QueryEscape(ctx.ContentID),
		&networking.RequestParams{
			Headers: youtubeWebHeaders(),
		},
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch youtube page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to get youtube page: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read youtube page: %w", err)
	}

	var apiKey string
	if match := innertubeKeyPattern.FindSubmatch(body); len(match) >= 2 {
		apiKey = string(match[1])
	}
	return body, apiKey, nil
}

func parseInitialPlayerResponse(body []byte) *PlayerResponse {
	match := initialPlayerRespPattern.FindSubmatch(body)
	if len(match) < 2 {
		return nil
	}

	var response PlayerResponse
	if err := sonic.ConfigFastest.Unmarshal(match[1], &response); err != nil {
		return nil
	}
	return &response
}

func fetchPlayerAPI(
	ctx *models.ExtractorContext,
	apiKey string,
	client playerClient,
) (*PlayerResponse, error) {
	endpoint := "https://www.youtube.com/youtubei/v1/player?prettyPrint=false"
	if apiKey != "" {
		endpoint += "&key=" + url.QueryEscape(apiKey)
	}

	clientPayload := map[string]any{
		"clientName":    client.Name,
		"clientVersion": client.Version,
	}
	if client.AndroidSDK != 0 {
		clientPayload["androidSdkVersion"] = client.AndroidSDK
	}

	payload := map[string]any{
		"context": map[string]any{
			"client": clientPayload,
		},
		"videoId":        ctx.ContentID,
		"contentCheckOk": true,
		"racyCheckOk":    true,
	}
	body, _ := sonic.ConfigFastest.Marshal(payload)

	headers := youtubeWebHeaders()
	headers["Accept"] = "application/json"
	headers["Content-Type"] = "application/json"
	headers["Origin"] = "https://www.youtube.com"
	headers["Referer"] = "https://www.youtube.com/watch?v=" + ctx.ContentID
	if client.UserAgent != "" {
		headers["User-Agent"] = client.UserAgent
	}

	resp, err := ctx.Fetch(
		http.MethodPost,
		endpoint,
		&networking.RequestParams{
			Body:    bytes.NewReader(body),
			Headers: headers,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch youtube player api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid youtube player api response: %s", resp.Status)
	}

	var playerResponse PlayerResponse
	if err := sonic.ConfigFastest.NewDecoder(resp.Body).Decode(&playerResponse); err != nil {
		return nil, fmt.Errorf("failed to parse youtube player api response: %w", err)
	}
	return &playerResponse, nil
}

func hasUsableFormats(resp *PlayerResponse) bool {
	if resp == nil {
		return false
	}
	for _, format := range append(resp.StreamingData.Formats, resp.StreamingData.AdaptiveFormats...) {
		if isVideoFormat(format) && directFormatURL(format) != "" {
			return true
		}
	}
	return false
}

func directFormatURL(format Format) string {
	if format.URL != "" {
		return util.UnescapeURL(format.URL)
	}
	cipher := format.SignatureCipher
	if cipher == "" {
		cipher = format.Cipher
	}
	if cipher == "" {
		return ""
	}

	values, err := url.ParseQuery(cipher)
	if err != nil {
		return ""
	}
	// Some clients return a cipher object that is already signed. If an "s"
	// signature is present it needs YouTube's JS decipher step, so skip it.
	if values.Get("s") != "" {
		return ""
	}
	return util.UnescapeURL(values.Get("url"))
}

func playabilityError(resp *PlayerResponse) error {
	if resp == nil {
		return fmt.Errorf("youtube player response not found")
	}
	status := resp.PlayabilityStatus.Status
	reason := resp.PlayabilityStatus.Reason
	switch status {
	case "LOGIN_REQUIRED":
		return util.ErrAuthenticationNeeded
	case "AGE_CHECK_REQUIRED":
		return util.ErrAgeRestricted
	case "UNPLAYABLE":
		if strings.Contains(strings.ToLower(reason), "age") {
			return util.ErrAgeRestricted
		}
		if strings.Contains(strings.ToLower(reason), "members-only") {
			return util.ErrPaidContent
		}
	}
	if reason != "" {
		return fmt.Errorf("youtube is not playable: %s", reason)
	}
	if status != "" && status != "OK" {
		return fmt.Errorf("youtube is not playable: %s", status)
	}
	return fmt.Errorf("no usable youtube formats found")
}

func parseInt32(value string) int32 {
	parsed, _ := strconv.ParseInt(value, 10, 32)
	return int32(parsed)
}

func parseInt64(value string) int64 {
	parsed, _ := strconv.ParseInt(value, 10, 64)
	return parsed
}

func formatDurationSeconds(videoDetailsSeconds string, format Format) int32 {
	if seconds := parseInt32(videoDetailsSeconds); seconds != 0 {
		return seconds
	}
	if format.ApproxDurationMS == "" {
		return 0
	}
	return int32(parseInt64(format.ApproxDurationMS) / 1000)
}

func youtubeWebHeaders() map[string]string {
	return map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Accept-Language": "en-US,en;q=0.9",
	}
}
