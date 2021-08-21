package avakian

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/holedaemon/avakian/internal/version"
	"github.com/skwair/harmony/discord"
)

const (
	discordAPIRoot = "https://discord.com/api/v9"
)

var (
	userAgent = "DiscordBot (https://github.com/holedaemon/avakian, " + version.Version() + ")"
)

type modifyCurrentUserRequest struct {
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

func (b *Bot) ModifyCurrentUser(ctx context.Context, username, avatar string) (*discord.User, error) {
	cu := &modifyCurrentUserRequest{
		Username: username,
		Avatar:   avatar,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(cu); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, discordAPIRoot+"/users/@me", buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", b.Token))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	res, err := b.HTTP.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, discord.NewAPIError(res)
	}

	var u *discord.User
	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		return nil, err
	}

	return u, nil
}
