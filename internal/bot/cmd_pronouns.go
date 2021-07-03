package bot

import (
	"context"
	"strings"

	"github.com/erei/avakian/internal/database/models"
	"github.com/erei/avakian/internal/pkg/modelsx"
	"github.com/erei/avakian/internal/pkg/zapx"
	"github.com/skwair/harmony/discord"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/zikaeroh/ctxlog"
)

var (
	cmdPronouns = &MessageCommand{
		permissions: discord.PermissionSendMessages,
		fn:          cmdPronounsFn,
	}

	pronounsCommands = map[string]*MessageCommand{
		"init": {
			permissions: discord.PermissionManageRoles,
			fn:          cmdPronounsInit,
		},
		"add": {
			permissions: discord.PermissionSendMessages,
			fn:          cmdPronounsAdd,
		},
		"remove": {
			permissions: discord.PermissionSendMessages,
			fn:          cmdPronounsRemove,
		},
		"create": {
			permissions: discord.PermissionManageRoles,
			fn:          cmdPronounsCreate,
		},
		"delete": {
			permissions: discord.PermissionManageRoles,
			fn:          cmdPronounsDelete,
		},
		"list": {
			permissions: discord.PermissionSendMessages,
			fn:          cmdPronounsList,
		},
	}

	defaultPronouns = []string{
		"he/him",
		"she/her",
		"they/them",
		"they/any",
		"they/he",
		"they/she",
	}
)

func cmdPronounsFn(ctx context.Context, s *MessageSession) error {
	usage := func() error {
		return s.Replyf(ctx, "Usage: `%s`", buildUsage(s.Prefix, "pronouns", pronounsCommands))
	}

	if len(s.Args) == 0 {
		return usage()
	}

	subSess := *s
	subSess.Args = s.Args[1:]
	sub := s.Args[0]

	cmd := pronounsCommands[sub]
	if cmd == nil {
		return usage()
	}

	p, err := s.Bot.FetchMemberPermissions(ctx, s.Msg.GuildID, s.Msg.ChannelID, s.Msg.Author.ID)
	if err != nil {
		return err
	}

	if !cmd.HasPermission(p) {
		ctxlog.Debug(ctx, "member lacks permission to run command", zapx.Member(s.Msg.Author.ID))
		return nil
	}

	if err := cmd.Execute(ctx, &subSess); err != nil {
		return err
	}

	return nil
}

func cmdPronounsAdd(ctx context.Context, s *MessageSession) error {
	if len(s.Args) == 0 {
		return s.Reply(ctx, "At least one argument is required")
	}

	prn := s.Args[0]
	dp, err := modelsx.GetPronoun(ctx, s.Tx, s.Msg.GuildID, prn)
	if err != nil {
		return err
	}

	if dp == nil {
		return s.Reply(ctx, "Those pronouns haven't been added to this guild")
	}

	if err := s.Bot.AddRole(ctx, s.Msg.GuildID, s.Msg.Author.ID, dp.RoleSnowflake, "Automated; pronoun role requested"); err != nil {
		return err
	}

	return s.Replyf(ctx, "Thou hath been branded %s", strings.ToUpper(prn))
}

func cmdPronounsRemove(ctx context.Context, s *MessageSession) error {
	if len(s.Args) == 0 {
		return s.Reply(ctx, "At least one argument is required")
	}

	prn := s.Args[0]
	dp, err := modelsx.GetPronoun(ctx, s.Tx, s.Msg.GuildID, prn)
	if err != nil {
		return err
	}

	if dp == nil {
		return s.Reply(ctx, "Those pronouns haven't been added to this guild")
	}

	if err := s.Bot.RemoveRole(ctx, s.Msg.GuildID, s.Msg.Author.ID, dp.RoleSnowflake, "Automated; pronoun removal requested"); err != nil {
		return err
	}

	return s.Replyf(ctx, "Thou no longer bears the brand of %s", dp.Pronoun)
}

func cmdPronounsCreate(ctx context.Context, s *MessageSession) error {
	if len(s.Args) == 0 {
		return s.Reply(ctx, "At least one argument is required")
	}

	prn := s.Args[0]
	exists, err := modelsx.PronounExists(ctx, s.Tx, s.Msg.GuildID, prn)
	if err != nil {
		return err
	}

	if exists {
		return s.Reply(ctx, "Those pronouns have already been created")
	}

	role, err := s.Bot.CreateRole(ctx, s.Msg.GuildID, "Automated; creation of pronoun role", discord.WithRoleName(prn))
	if err != nil {
		return err
	}

	newPrn := &models.Pronoun{
		GuildSnowflake: s.Msg.GuildID,
		Pronoun:        strings.ToLower(prn),
		RoleSnowflake:  role.ID,
	}

	if err := newPrn.Insert(ctx, s.Tx, boil.Infer()); err != nil {
		return err
	}

	return s.Reply(ctx, "Pronouns added to guild")
}

func cmdPronounsDelete(ctx context.Context, s *MessageSession) error {
	if len(s.Args) == 0 {
		return s.Reply(ctx, "At least one argument is required")
	}

	prn := s.Args[0]
	id, deleted, err := modelsx.DeletePronoun(ctx, s.Tx, s.Msg.GuildID, prn)
	if err != nil {
		return err
	}

	if !deleted {
		return s.Reply(ctx, "Those pronouns don't exist on this guild")
	}

	if err := s.Bot.DeleteRole(ctx, s.Msg.GuildID, id, "Automated; pronoun removal requested"); err != nil {
		return err
	}

	return s.Reply(ctx, "Pronouns removed from guild")
}

func cmdPronounsList(ctx context.Context, s *MessageSession) error {
	list, err := modelsx.Pronouns(ctx, s.Tx, s.Msg.GuildID)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return s.Reply(ctx, "No pronouns have been added to this guild")
	}

	var sb strings.Builder

	sb.WriteString("```")

	for _, p := range list {
		sb.WriteString(p.Pronoun + "\n")
	}

	sb.WriteString("```")

	return s.Reply(ctx, sb.String())
}

func cmdPronounsInit(ctx context.Context, s *MessageSession) error {
	g, err := s.Bot.FetchGuild(ctx, s.Msg.GuildID)
	if err != nil {
		return err
	}

	needed := make(map[string]bool)
	for _, dp := range defaultPronouns {
		needed[dp] = true
	}

	for _, r := range g.Roles {
		_, ok := needed[r.Name]
		if ok {
			needed[r.Name] = false
		}
	}

	var sb strings.Builder
	errors := 0
	created := 0
	fmtError := func(name, kind string, err error) {
		if sb.Len() == 0 {
			sb.WriteString("Errors encountered during initialization:\n```")
		}

		sb.WriteString(name + "(" + kind + "):\t" + err.Error())
		errors++
	}

	for key, val := range needed {
		key = strings.ToLower(key)

		if val {
			r, err := s.Bot.CreateRole(ctx, s.Msg.GuildID, "Initializing pronoun roles...", discord.WithRoleName(key))
			if err != nil {
				fmtError(key, "api", err)
				continue
			}

			prn := &models.Pronoun{
				GuildSnowflake: s.Msg.GuildID,
				Pronoun:        key,
				RoleSnowflake:  r.ID,
			}

			if err := prn.Insert(ctx, s.Tx, boil.Infer()); err != nil {
				fmtError(key, "db", err)
				continue
			}

			created++
		}
	}

	if errors > 0 {
		sb.WriteString("```")
		return s.Reply(ctx, sb.String())
	}

	if created == 0 {
		return s.Reply(ctx, "0 pronouns initialized as they already existed")
	}

	return s.Reply(ctx, "Pronouns initialized successfully")
}
