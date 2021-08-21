package avakian

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/holedaemon/avakian/internal/bot/message"
	"github.com/holedaemon/avakian/internal/pkg/httpx"
	"github.com/zikaeroh/ctxlog"
)

const (
	discordNameFloor   = 2
	discordNameCeiling = 32
)

var (
	cmdAdminRename = message.NewCommand(
		message.WithCommandArgs(1),
		message.WithCommandFn(cmdAdminRenameFn),
		message.WithCommandPermissions(-1),
	)

	cmdAdminAvatar = message.NewCommand(
		message.WithCommandFn(cmdAdminAvatarFn),
		message.WithCommandPermissions(-1),
	)

	cmdAdmin = message.NewCommand(
		message.WithCommandFn(cmdAdminFn),
		message.WithCommandPermissions(-1),
	)

	adminCommands = message.NewCommandMap(
		message.WithMapScope(true),
		message.WithMapCommand("rename", cmdAdminRename),
		message.WithMapCommand("identitycrisis", cmdAdminRename),
		message.WithMapCommand("setavatar", cmdAdminAvatar),
		message.WithMapCommand("glowup", cmdAdminAvatar),
	)
)

func cmdAdminFn(ctx context.Context, s *message.Session) error {
	return adminCommands.ExecuteCommand(ctx, s)
}

func cmdAdminRenameFn(ctx context.Context, s *message.Session) error {
	newName := strings.Join(s.Args, " ")
	if newName == "" {
		return s.Reply(ctx, "Name's empty, dumbass")
	}

	ctxlog.Debug(ctx, newName)

	b := getBot(s)
	me := b.Client.State.Me()
	if me.Username == newName {
		return s.Reply(ctx, "New name is the same as my current")
	}

	if len(newName) < discordNameFloor || len(newName) > discordNameCeiling {
		ctxlog.Debug(ctx, "hello")
		return s.Replyf(ctx, "Username must be between %d and %d characaters in length", discordNameFloor, discordNameCeiling)
	}

	_, err := b.ModifyCurrentUser(ctx, newName, "")
	if err != nil {
		return err
	}

	return s.Reply(ctx, "Username has been updated")
}

func cmdAdminAvatarFn(ctx context.Context, s *message.Session) error {
	var avatar string
	if len(s.Args) > 0 {
		avatar = s.Args[0]
	} else {
		if len(s.Msg.Attachments) > 0 {
			avatar = s.Msg.Attachments[0].URL
		} else {
			return s.Reply(ctx, "New avatar wasn't provided")
		}
	}

	b := getBot(s)
	res, err := b.HTTP.Get(avatar)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("%w: %d", httpx.ErrStatusCode, res.StatusCode)
	}

	ct := res.Header.Get("Content-Type")
	switch ct {
	case "image/jpeg", "image/gif", "image/png":
	default:
		return s.Reply(ctx, "Bad Content-Type, homeslice")
	}

	buf := new(bytes.Buffer)
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("data:%s;base64,", ct))

	io.Copy(buf, res.Body)

	sb.WriteString(base64.RawStdEncoding.EncodeToString(buf.Bytes()))

	_, err = b.ModifyCurrentUser(ctx, "", sb.String())
	if err != nil {
		return err
	}

	return s.Reply(ctx, "Avatar has been updated")
}
