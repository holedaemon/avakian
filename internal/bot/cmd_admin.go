package bot

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
)

var (
	cmdAdmin = &MessageCommand{
		permissions: -1,
		fn:          cmdAdminFn,
	}

	adminCommands = map[string]*MessageCommand{
		"modify": {
			permissions: -1,
			fn:          cmdAdminModify,
		},
	}
)

func cmdAdminFn(ctx context.Context, s *MessageSession) error {
	return nil
}

func cmdAdminModify(ctx context.Context, s *MessageSession) error {
	if len(s.Args) == 0 {
		return s.Reply(ctx, "If you want me to go through an identity crisis you're gonna have to supply a new name or avatar.")
	}

	newName := s.Args[0]
	newAvatar := s.Args[1]

	user := s.Bot.Client.User("@me")

	var sb strings.Builder
	sb.WriteString("data:")

	if newAvatar == "" {
		me, err := user.Get(ctx)
		if err != nil {
			return err
		}

		newAvatar = me.AvatarURL()
	}

	res, err := http.Get(newAvatar)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	buf := new(bytes.Buffer)
	written, err := io.Copy(buf, res.Body)
	if err != nil && written == 0 {
		return err
	}

	mimeType := http.DetectContentType(buf.Bytes())

	switch mimeType {
	case "image/jpeg":
		sb.WriteString("image/jpeg")
	case "image/png":
		sb.WriteString("image/png")
	case "image/gif":
		sb.WriteString("image/gif")
	default:
		return s.Replyf(ctx, "Invalid image mime type: %s", mimeType)
	}

	sb.WriteString(";base64,")

	sb.WriteString(
		base64.RawURLEncoding.EncodeToString(buf.Bytes()),
	)

	_, err = user.Modify(ctx, newName, sb.String())
	if err != nil {
		return err
	}

	return s.Reply(ctx, "My name and/or avatar has been updated successfully")
}
