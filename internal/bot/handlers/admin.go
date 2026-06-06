package handlers

import (
	"context"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"eadownloader/internal/config"
	"eadownloader/internal/database"
	"eadownloader/internal/util"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	adminCallbackPrefix = "admin:"

	adminScreenHome       = "home"
	adminScreenModeration = "moderation"
	adminScreenSystem     = "system"
	adminScreenUsers      = "users"
	adminScreenBans       = "bans"
	adminScreenUser       = "user"

	adminActionBanConfirm = "ban_confirm"
	adminActionBan        = "ban"
	adminActionUnban      = "unban"
)

func AdminHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	if !util.IsBotAdmin(ctx) {
		return ext.EndGroups
	}

	text, keyboard, err := buildAdminHome()
	if err != nil {
		return err
	}

	ctx.EffectiveMessage.Reply(
		bot,
		text,
		&gotgbot.SendMessageOpts{
			ParseMode:   gotgbot.ParseModeHTML,
			ReplyMarkup: keyboard,
		},
	)
	return ext.EndGroups
}

func AdminCallbackHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.CallbackQuery == nil || !util.IsAdminID(ctx.CallbackQuery.From.Id) {
		return ext.EndGroups
	}

	text, keyboard, err := resolveAdminCallback(ctx)
	if err != nil {
		return err
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

func resolveAdminCallback(ctx *ext.Context) (string, gotgbot.InlineKeyboardMarkup, error) {
	data := strings.TrimPrefix(ctx.CallbackQuery.Data, adminCallbackPrefix)

	switch {
	case data == adminScreenHome:
		return buildAdminHome()
	case data == adminScreenModeration:
		return buildModerationHome()
	case data == adminScreenUsers:
		return buildUserList()
	case data == adminScreenBans:
		return buildBannedUserList()
	case data == adminScreenSystem:
		return buildSystemPanel()
	case strings.HasPrefix(data, adminScreenUser+":"):
		return buildUserProfile(strings.TrimPrefix(data, adminScreenUser+":"))
	case strings.HasPrefix(data, adminActionBanConfirm+":"):
		return buildBanConfirm(strings.TrimPrefix(data, adminActionBanConfirm+":"))
	case strings.HasPrefix(data, adminActionBan+":"):
		return banUserFromCallback(ctx, strings.TrimPrefix(data, adminActionBan+":"))
	case strings.HasPrefix(data, adminActionUnban+":"):
		return unbanUserFromCallback(strings.TrimPrefix(data, adminActionUnban+":"))
	default:
		return buildAdminHome()
	}
}

func buildAdminHome() (string, gotgbot.InlineKeyboardMarkup, error) {
	stats, err := database.Q().GetStats(
		context.Background(),
		pgtype.Timestamptz{
			Time:  time.Now().Add(-100 * 365 * 24 * time.Hour),
			Valid: true,
		},
	)
	if err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}
	bannedUsersCount, err := database.Q().CountBannedUsers(context.Background())
	if err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}

	text := fmt.Sprintf(
		"<b>EaDownloader Admin</b>\n\n"+
			"<b>Overview</b>\n"+
			"Users: %d\n"+
			"Groups: %d\n"+
			"Downloads: %d\n"+
			"Storage: %s\n"+
			"Banned users: %d\n\n"+
			"Choose a module below.",
		stats.TotalPrivateChats,
		stats.TotalGroupChats,
		stats.TotalDownloads,
		formatBytes(stats.TotalDownloadsSize),
		bannedUsersCount,
	)

	return text, gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Analytics", CallbackData: statsCallbackPrefix + statsScreenSummary + ":" + statsPeriodAll},
				{Text: "Moderation", CallbackData: adminCallbackPrefix + adminScreenModeration},
			},
			{
				{Text: "System", CallbackData: adminCallbackPrefix + adminScreenSystem},
				{Text: "Close", CallbackData: "close"},
			},
		},
	}, nil
}

func buildModerationHome() (string, gotgbot.InlineKeyboardMarkup, error) {
	bannedUsersCount, err := database.Q().CountBannedUsers(context.Background())
	if err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}

	text := fmt.Sprintf(
		"<b>Moderation</b>\n\n"+
			"Manage users from a clean profile view.\n"+
			"Banned users: %d",
		bannedUsersCount,
	)

	return text, gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Users", CallbackData: adminCallbackPrefix + adminScreenUsers},
				{Text: "Banned Users", CallbackData: adminCallbackPrefix + adminScreenBans},
			},
			{
				{Text: "Back", CallbackData: adminCallbackPrefix + adminScreenHome},
				{Text: "Close", CallbackData: "close"},
			},
		},
	}, nil
}

func buildUserList() (string, gotgbot.InlineKeyboardMarkup, error) {
	rows, err := database.Q().ListChatsByType(
		context.Background(),
		database.ListChatsByTypeParams{
			Type:       database.ChatTypePrivate,
			LimitCount: statsListLimit,
		},
	)
	if err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}

	if len(rows) == 0 {
		return "<b>Users</b>\n\nNo users recorded yet.", userListKeyboard(rows), nil
	}

	text := fmt.Sprintf("<b>Users</b>\nLast %d active private users\n\n", len(rows))
	for index, row := range rows {
		status := "active"
		if banned, err := database.Q().IsUserBanned(context.Background(), row.ChatID); err == nil && banned {
			status = "banned"
		}
		text += fmt.Sprintf(
			"<b>%d.</b> %s\n<code>%d</code> · %s · %s\n\n",
			index+1,
			formatChatDisplayName(row),
			row.ChatID,
			status,
			formatTimeAgo(row.LastSeenAt),
		)
	}

	return strings.TrimSpace(text), userListKeyboard(rows), nil
}

func buildBannedUserList() (string, gotgbot.InlineKeyboardMarkup, error) {
	total, err := database.Q().CountBannedUsers(context.Background())
	if err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}
	rows, err := database.Q().ListBannedUsers(context.Background(), statsListLimit)
	if err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}

	if len(rows) == 0 {
		return "<b>Banned Users</b>\n\nNo banned users yet.", bannedUserListKeyboard(rows), nil
	}

	text := fmt.Sprintf("<b>Banned Users</b>\nTotal: %d\n\n", total)
	for index, row := range rows {
		text += fmt.Sprintf(
			"<b>%d.</b> %s\n<code>%d</code> · %s\nReason: %s\n\n",
			index+1,
			formatBannedUserDisplayName(row),
			row.UserID,
			formatTimeAgo(row.CreatedAt),
			html.EscapeString(row.Reason),
		)
	}

	return strings.TrimSpace(text), bannedUserListKeyboard(rows), nil
}

func buildUserProfile(value string) (string, gotgbot.InlineKeyboardMarkup, error) {
	userID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return buildUserList()
	}

	user, err := database.Q().GetChatByID(context.Background(), userID)
	if err != nil {
		return buildUnknownUserProfile(userID)
	}

	banned, err := database.Q().IsUserBanned(context.Background(), user.ChatID)
	if err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}

	status := "active"
	if banned {
		status = "banned"
	}

	text := fmt.Sprintf(
		"<b>User Profile</b>\n\n"+
			"%s\n"+
			"ID: <code>%d</code>\n"+
			"Username: %s\n"+
			"Language: %s\n"+
			"Status: %s\n"+
			"Created: %s\n"+
			"Last seen: %s",
		formatUserProfileDisplayName(user),
		user.ChatID,
		formatUsername(user.Username),
		html.EscapeString(user.Language),
		status,
		formatTimeAgo(user.CreatedAt),
		formatTimeAgo(user.LastSeenAt),
	)

	return text, userProfileKeyboard(user.ChatID, banned), nil
}

func buildUnknownUserProfile(userID int64) (string, gotgbot.InlineKeyboardMarkup, error) {
	banned, err := database.Q().IsUserBanned(context.Background(), userID)
	if err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}

	text := fmt.Sprintf(
		"<b>User Profile</b>\n\n"+
			"ID: <code>%d</code>\n"+
			"Status: %s\n\n"+
			"This user is not in the chat table yet.",
		userID,
		map[bool]string{true: "banned", false: "unknown"}[banned],
	)
	return text, userProfileKeyboard(userID, banned), nil
}

func buildBanConfirm(value string) (string, gotgbot.InlineKeyboardMarkup, error) {
	userID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return buildUserList()
	}
	if util.IsAdminID(userID) {
		return "<b>Protected User</b>\n\nAdmins cannot be banned.", userProfileKeyboard(userID, false), nil
	}

	text := fmt.Sprintf(
		"<b>Confirm Ban</b>\n\n"+
			"User ID: <code>%d</code>\n\n"+
			"The user will be blocked from private chats, groups and inline mode.",
		userID,
	)
	return text, gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Confirm Ban", CallbackData: adminCallbackPrefix + adminActionBan + ":" + strconv.FormatInt(userID, 10)},
			},
			{
				{Text: "Back", CallbackData: adminCallbackPrefix + adminScreenUser + ":" + strconv.FormatInt(userID, 10)},
				{Text: "Close", CallbackData: "close"},
			},
		},
	}, nil
}

func buildSystemPanel() (string, gotgbot.InlineKeyboardMarkup, error) {
	bannedUsersCount, err := database.Q().CountBannedUsers(context.Background())
	if err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}

	text := fmt.Sprintf(
		"<b>System</b>\n\n"+
			"Admins: %d\n"+
			"Whitelist: %d\n"+
			"Banned users: %d\n"+
			"Concurrent updates: %d\n"+
			"Max duration: %s\n"+
			"Max file size: %s\n"+
			"Caching: %t\n"+
			"Log level: %s\n"+
			"Time: %s",
		len(config.Env.Admins),
		len(config.Env.Whitelist),
		bannedUsersCount,
		config.Env.ConcurrentUpdates,
		config.Env.MaxDuration,
		formatBytes(config.Env.MaxFileSize),
		config.Env.Caching,
		config.Env.LogLevel.String(),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return text, adminBackKeyboard(adminScreenHome), nil
}

func banUserFromCallback(ctx *ext.Context, value string) (string, gotgbot.InlineKeyboardMarkup, error) {
	userID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return buildUserList()
	}
	if util.IsAdminID(userID) {
		return "<b>Protected User</b>\n\nAdmins cannot be banned.", userProfileKeyboard(userID, false), nil
	}

	_, err = database.Q().BanUser(
		context.Background(),
		database.BanUserParams{
			UserID:   userID,
			Reason:   "admin panel",
			BannedBy: ctx.CallbackQuery.From.Id,
		},
	)
	if err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}
	return buildUserProfile(value)
}

func unbanUserFromCallback(value string) (string, gotgbot.InlineKeyboardMarkup, error) {
	userID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return buildUserList()
	}
	if err := database.Q().UnbanUser(context.Background(), userID); err != nil {
		return "", gotgbot.InlineKeyboardMarkup{}, err
	}
	return buildUserProfile(value)
}

func userListKeyboard(rows []database.ListChatsByTypeRow) gotgbot.InlineKeyboardMarkup {
	buttons := numberedUserButtons(rows)
	buttons = append(buttons, []gotgbot.InlineKeyboardButton{
		{Text: "Banned Users", CallbackData: adminCallbackPrefix + adminScreenBans},
	})
	buttons = append(buttons, adminBackRow(adminScreenModeration))
	return gotgbot.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func bannedUserListKeyboard(rows []database.ListBannedUsersRow) gotgbot.InlineKeyboardMarkup {
	buttons := make([][]gotgbot.InlineKeyboardButton, 0, 4)
	numberedButtons := make([]gotgbot.InlineKeyboardButton, 0, len(rows))
	for index, row := range rows {
		numberedButtons = append(numberedButtons, gotgbot.InlineKeyboardButton{
			Text:         strconv.Itoa(index + 1),
			CallbackData: adminCallbackPrefix + adminScreenUser + ":" + strconv.FormatInt(row.UserID, 10),
		})
	}
	buttons = append(buttons, chunkButtons(numberedButtons, 5)...)
	buttons = append(buttons, []gotgbot.InlineKeyboardButton{
		{Text: "Users", CallbackData: adminCallbackPrefix + adminScreenUsers},
	})
	buttons = append(buttons, adminBackRow(adminScreenModeration))
	return gotgbot.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func numberedUserButtons(rows []database.ListChatsByTypeRow) [][]gotgbot.InlineKeyboardButton {
	buttons := make([]gotgbot.InlineKeyboardButton, 0, len(rows))
	for index, row := range rows {
		buttons = append(buttons, gotgbot.InlineKeyboardButton{
			Text:         strconv.Itoa(index + 1),
			CallbackData: adminCallbackPrefix + adminScreenUser + ":" + strconv.FormatInt(row.ChatID, 10),
		})
	}
	return chunkButtons(buttons, 5)
}

func userProfileKeyboard(userID int64, banned bool) gotgbot.InlineKeyboardMarkup {
	actionText := "Ban User"
	actionData := adminCallbackPrefix + adminActionBanConfirm + ":" + strconv.FormatInt(userID, 10)
	if banned {
		actionText = "Unban User"
		actionData = adminCallbackPrefix + adminActionUnban + ":" + strconv.FormatInt(userID, 10)
	}

	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: actionText, CallbackData: actionData},
			},
			{
				{Text: "Users", CallbackData: adminCallbackPrefix + adminScreenUsers},
				{Text: "Banned Users", CallbackData: adminCallbackPrefix + adminScreenBans},
			},
			adminBackRow(adminScreenModeration),
		},
	}
}

func adminBackKeyboard(screen string) gotgbot.InlineKeyboardMarkup {
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			adminBackRow(screen),
		},
	}
}

func adminBackRow(screen string) []gotgbot.InlineKeyboardButton {
	return []gotgbot.InlineKeyboardButton{
		{Text: "Back", CallbackData: adminCallbackPrefix + screen},
		{Text: "Close", CallbackData: "close"},
	}
}

func chunkButtons(buttons []gotgbot.InlineKeyboardButton, size int) [][]gotgbot.InlineKeyboardButton {
	rows := make([][]gotgbot.InlineKeyboardButton, 0, (len(buttons)+size-1)/size)
	for start := 0; start < len(buttons); start += size {
		end := min(start+size, len(buttons))
		rows = append(rows, buttons[start:end])
	}
	return rows
}

func formatBannedUserDisplayName(row database.ListBannedUsersRow) string {
	name := bannedUserDisplayLabel(row)
	return fmt.Sprintf(
		"<a href='tg://user?id=%d'>%s</a>",
		row.UserID,
		html.EscapeString(name),
	)
}

func bannedUserDisplayLabel(row database.ListBannedUsersRow) string {
	name := strings.TrimSpace(joinValidTexts(row.FirstName, row.LastName))
	if name == "" && row.Username.Valid && row.Username.String != "" {
		name = "@" + row.Username.String
	}
	if name == "" {
		name = strconv.FormatInt(row.UserID, 10)
	}
	return name
}

func joinValidTexts(values ...pgtype.Text) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		if value.Valid && strings.TrimSpace(value.String) != "" {
			parts = append(parts, strings.TrimSpace(value.String))
		}
	}
	return strings.Join(parts, " ")
}

func formatUsername(username string) string {
	if strings.TrimSpace(username) == "" {
		return "-"
	}
	return "@" + html.EscapeString(username)
}

func formatUserProfileDisplayName(user database.GetChatByIDRow) string {
	name := strings.TrimSpace(user.Title)
	if name == "" {
		name = strings.TrimSpace(strings.Join([]string{user.FirstName, user.LastName}, " "))
	}
	if name == "" && user.Username != "" {
		name = "@" + user.Username
	}
	if name == "" {
		name = strconv.FormatInt(user.ChatID, 10)
	}
	if user.Type == database.ChatTypePrivate {
		return fmt.Sprintf(
			"<a href='tg://user?id=%d'>%s</a>",
			user.ChatID,
			html.EscapeString(name),
		)
	}
	return html.EscapeString(name)
}
