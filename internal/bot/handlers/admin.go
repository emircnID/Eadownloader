package handlers

import (
	"fmt"
	"time"

	"eadownloader/internal/config"
	"eadownloader/internal/util"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const (
	adminCallbackPrefix = "admin:"

	adminScreenHome   = "home"
	adminScreenSystem = "system"
)

func AdminHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	if !util.IsBotAdmin(ctx) {
		return ext.EndGroups
	}

	ctx.EffectiveMessage.Reply(
		bot,
		formatAdminHome(),
		&gotgbot.SendMessageOpts{
			ParseMode:   gotgbot.ParseModeHTML,
			ReplyMarkup: getAdminKeyboard(),
		},
	)
	return ext.EndGroups
}

func AdminCallbackHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.CallbackQuery == nil || !util.IsAdminID(ctx.CallbackQuery.From.Id) {
		return ext.EndGroups
	}

	text := formatAdminHome()
	keyboard := getAdminKeyboard()
	if ctx.CallbackQuery.Data == adminCallbackPrefix+adminScreenSystem {
		text = formatAdminSystem()
		keyboard = getAdminSystemKeyboard()
	}

	ctx.CallbackQuery.Answer(bot, nil)
	ctx.EffectiveMessage.EditText(
		bot,
		text,
		&gotgbot.EditMessageTextOpts{
			ParseMode:   gotgbot.ParseModeHTML,
			ReplyMarkup: keyboard,
		},
	)
	return nil
}

func formatAdminHome() string {
	return "<b>EaDownloader admin</b>\n\n" +
		"Stats, errors and quick controls are collected here."
}

func formatAdminSystem() string {
	return fmt.Sprintf(
		"<b>System</b>\n\n"+
			"Admins: %d\n"+
			"Whitelist: %d\n"+
			"Concurrent updates: %d\n"+
			"Max duration: %s\n"+
			"Max file size: %s\n"+
			"Caching: %t\n"+
			"Log level: %s\n"+
			"Time: %s",
		len(config.Env.Admins),
		len(config.Env.Whitelist),
		config.Env.ConcurrentUpdates,
		config.Env.MaxDuration,
		formatBytes(config.Env.MaxFileSize),
		config.Env.Caching,
		config.Env.LogLevel.String(),
		time.Now().Format("2006-01-02 15:04:05"),
	)
}

func getAdminKeyboard() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Stats", CallbackData: statsCallbackPrefix + statsScreenSummary + ":" + statsPeriodAll},
				{Text: "Errors", CallbackData: statsCallbackPrefix + statsScreenErrors},
			},
			{
				{Text: "Platforms", CallbackData: statsCallbackPrefix + statsScreenPlatforms + ":" + statsPeriodAll},
				{Text: "System", CallbackData: adminCallbackPrefix + adminScreenSystem},
			},
			{
				{Text: "Users", CallbackData: statsCallbackPrefix + statsScreenUsers},
				{Text: "Groups", CallbackData: statsCallbackPrefix + statsScreenGroups},
			},
			{
				{Text: "Close", CallbackData: "close"},
			},
		},
	}
}

func getAdminSystemKeyboard() gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Back", CallbackData: adminCallbackPrefix + adminScreenHome},
				{Text: "Close", CallbackData: "close"},
			},
		},
	}
}
