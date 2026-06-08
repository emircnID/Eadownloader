package handlers

import (
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func TestTextMentionModerationTarget(t *testing.T) {
	user := &gotgbot.User{
		Id:        6473764544,
		FirstName: "Eclaire",
		LastName:  "Nuage",
		Username:  "EclaireNuage",
	}
	ctx := &ext.Context{
		EffectiveMessage: &gotgbot.Message{
			Text: "/ban Eclaire Nuage",
			Entities: []gotgbot.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 4},
				{Type: "text_mention", Offset: 5, Length: 13, User: user},
			},
		},
	}

	got, ok := textMentionModerationTarget(ctx)
	if !ok {
		t.Fatal("expected text mention target")
	}
	if got.Id != user.Id {
		t.Fatalf("expected user ID %d, got %d", user.Id, got.Id)
	}
}
