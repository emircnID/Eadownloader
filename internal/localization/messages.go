package localization

import "github.com/nicksnyder/go-i18n/v2/i18n"

var (
	Language = &i18n.Message{
		ID:    "Language",
		Other: "English",
	}
	StartMessage = &i18n.Message{
		ID:    "StartMessage",
		Other: "Welcome {{.Name}}. Send a supported social media link, and I will fetch the media for you.",
	}
	AddButton = &i18n.Message{
		ID:    "AddButton",
		Other: "Add to a Group",
	}
	ErrorMessage = &i18n.Message{
		ID:    "ErrorMessage",
		Other: "An error occurred. Please try again later.",
	}
	AddedToGroupMessage = &i18n.Message{
		ID:    "AddedToGroupMessage",
		Other: "Thank you for adding me! Use the /settings command to configure the bot for this group.",
	}
	SettingsButton = &i18n.Message{
		ID:    "SettingsButton",
		Other: "Settings",
	}
	LanguageButton = &i18n.Message{
		ID:    "LanguageButton",
		Other: "Language",
	}
	PrivateSettingsMessage = &i18n.Message{
		ID:    "PrivateSettingsMessage",
		Other: "Use the buttons below to change your personal bot settings.",
	}
	GroupSettingsMessage = &i18n.Message{
		ID:    "GroupSettingsMessage",
		Other: "Use the buttons below to change this group's bot settings.",
	}
	BackButton = &i18n.Message{
		ID:    "BackButton",
		Other: "Back",
	}
	SelectLanguageMessage = &i18n.Message{
		ID:    "SelectLanguageMessage",
		Other: "Select your preferred language.",
	}
	CaptionsSettingsMessage = &i18n.Message{
		ID:    "CaptionsSettingsMessage",
		Other: "When enabled, adds the original description to downloaded content if available.",
	}
	NsfwSettingsMessage = &i18n.Message{
		ID:    "NsfwSettingsMessage",
		Other: "When enabled, allows NSFW content to be downloaded in this chat.\n\nWarning: This type of content may violate Telegram's Terms of Service and result in group restrictions.",
	}
	SilentModeSettingsMessage = &i18n.Message{
		ID:    "SilentModeSettingsMessage",
		Other: "When enabled, the bot will not send error messages.",
	}
	MediaAlbumSettingsMessage = &i18n.Message{
		ID:    "MediaAlbumSettingsMessage",
		Other: "Select the maximum number of files allowed in a single media album.",
	}
	InlineLoadingMessage = &i18n.Message{
		ID:    "InlineLoadingMessage",
		Other: "Loading... Please wait.",
	}
	InlineProcessingMessage = &i18n.Message{
		ID:    "InlineProcessingMessage",
		Other: "Media shared! Processing the download... Please wait.",
	}
	InlineShareMessage = &i18n.Message{
		ID:    "InlineShareMessage",
		Other: "Share this media",
	}
	NoPermission = &i18n.Message{
		ID:    "NoPermission",
		Other: "You do not have permission to perform this action.",
	}
	CloseButton = &i18n.Message{
		ID:    "CloseButton",
		Other: "Close",
	}
	MediaAlbumButton = &i18n.Message{
		ID:    "MediaAlbumButton",
		Other: "Media Album",
	}
	SilentModeButton = &i18n.Message{
		ID:    "SilentModeButton",
		Other: "Silent Mode",
	}
	CaptionsButton = &i18n.Message{
		ID:    "CaptionsButton",
		Other: "Captions",
	}
	NsfwButton = &i18n.Message{
		ID:    "NsfwButton",
		Other: "NSFW",
	}
	ExtractorsButton = &i18n.Message{
		ID:    "ExtractorsButton",
		Other: "Platforms",
	}
	DisabledExtractorsSettingsMessage = &i18n.Message{
		ID:    "DisabledExtractorsSettingsMessage",
		Other: "Select which platforms should be disabled. Links from disabled platforms will be ignored by the bot.",
	}
	DeleteLinksButton = &i18n.Message{
		ID:    "DeleteProcessedButton",
		Other: "Links",
	}
	DeleteLinksSettingsMessage = &i18n.Message{
		ID:    "DeleteProcessedSettingsMessage",
		Other: "When enabled, deletes the user's original message after the link is processed successfully.",
	}
	SupportedExtractorsMessage = &i18n.Message{
		ID:    "SupportedExtractorsMessage",
		Other: "Supported platforms",
	}
	EnabledButton = &i18n.Message{
		ID:    "EnabledButton",
		Other: "Enabled",
	}
	DisabledButton = &i18n.Message{
		ID:    "DisabledButton",
		Other: "Disabled",
	}
	ErrorUnavailable = &i18n.Message{
		ID:    "ErrorUnavailable",
		Other: "This content is unavailable.",
	}
	ErrorTimeout = &i18n.Message{
		ID:    "ErrorTimeout",
		Other: "The download timed out. Please try again later.",
	}
	ErrorUnsupportedImageFormat = &i18n.Message{
		ID:    "ErrorUnsupportedImageFormat",
		Other: "Unsupported image format.",
	}
	ErrorUnsupportedExtractorType = &i18n.Message{
		ID:    "ErrorUnsupportedExtractorType",
		Other: "Unsupported platform type.",
	}
	ErrorMediaAlbumLimitExceeded = &i18n.Message{
		ID:    "ErrorMediaAlbumLimitExceeded",
		Other: "The media album limit exceeds the maximum allowed for this group. Use /settings to increase the limit.",
	}
	ErrorMediaAlbumGlobalLimitExceeded = &i18n.Message{
		ID:    "ErrorMediaAlbumGlobalLimitExceeded",
		Other: "The media album limit exceeds the maximum allowed for this instance.",
	}
	ErrorGeoRestrictedContent = &i18n.Message{
		ID:    "ErrorGeoRestrictedContent",
		Other: "This content has geo-restrictions and cannot be accessed from the server's location.",
	}
	ErrorNSFWNotAllowed = &i18n.Message{
		ID:    "ErrorNSFWNotAllowed",
		Other: "This content is marked as NSFW and cannot be downloaded in this group. Use /settings to allow NSFW content or use the bot privately.",
	}
	ErrorInlineMediaAlbum = &i18n.Message{
		ID:    "ErrorInlineMediaAlbum",
		Other: "Media albums cannot be downloaded in inline mode. Use the bot in a group or private chat.",
	}
	ErrorAuthenticationNeeded = &i18n.Message{
		ID:    "ErrorAuthenticationNeeded",
		Other: "This instance is not authenticated with this service.",
	}
	ErrorFileTooLarge = &i18n.Message{
		ID:    "ErrorFileTooLarge",
		Other: "This file is too large and exceeds the maximum allowed size for this instance.",
	}
	ErrorTelegramFileTooLarge = &i18n.Message{
		ID:    "ErrorTelegramFileTooLarge",
		Other: "This file is too large for Telegram.",
	}
	ErrorDurationTooLong = &i18n.Message{
		ID:    "ErrorDurationTooLong",
		Other: "This video is too long and exceeds the maximum allowed duration for this instance.",
	}
	ErrorPaidContent = &i18n.Message{
		ID:    "ErrorPaidContent",
		Other: "This content is paid and requires a subscription to access.",
	}
	ErrorAgeRestricted = &i18n.Message{
		ID:    "ErrorAgeRestricted",
		Other: "This content is age-restricted and cannot be accessed.",
	}
	ErrorPermissionDenied = &i18n.Message{
		ID:    "ErrorPermissionDenied",
		Other: "The bot does not have sufficient permissions to send this media. Please grant the necessary permissions and try again.",
	}
)
