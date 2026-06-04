package handlers

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"eadownloader/internal/localization"
	"eadownloader/internal/util"
)

// prevents the bot from processing a large
// backlog of messages after connection interruptions.
func OldMessagesHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	ts := ctx.EffectiveMessage.Date
	if time.Since(time.Unix(ts, 0)) > 2*time.Minute {
		return ext.EndGroups
	}
	return ext.ContinueGroups
}

func CloseHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	chat, err := util.ChatFromContext(ctx)
	if err != nil {
		return err
	}
	localizer := localization.New(chat.Language)
	isAdmin := util.CheckAdminPermission(bot, ctx, localizer)
	if !isAdmin {
		return nil
	}
	ctx.CallbackQuery.Answer(bot, nil)
	ctx.EffectiveMessage.Delete(bot, nil)
	return nil
}
