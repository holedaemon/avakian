package avakian

import (
	"context"
	"strings"

	"github.com/erei/avakian/internal/bot/message"
	"github.com/erei/avakian/internal/database/models"
	"github.com/erei/avakian/internal/pkg/modelsx"
	"github.com/erei/avakian/internal/pkg/snowflake"
	"github.com/skwair/harmony/discord"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	cmdPronouns = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionSendMessages),
		message.WithCommandUsage(pronounsCommands.BuildUsage("pronouns")),
		message.WithCommandFn(cmdPronounsFn),
	)

	cmdPronounsInit = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionManageRoles),
		message.WithCommandFn(cmdPronounsInitFn),
	)

	cmdPronounsAdd = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionSendMessages),
		message.WithCommandFn(cmdPronounsAddFn),
	)

	cmdPronounsRemove = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionSendMessages),
		message.WithCommandFn(cmdPronounsRemoveFn),
	)

	cmdPronounsCreate = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionManageRoles),
		message.WithCommandFn(cmdPronounsCreateFn),
	)

	cmdPronounsDelete = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionManageRoles),
		message.WithCommandFn(cmdPronounsDeleteFn),
	)

	cmdPronounsImport = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionManageRoles),
		message.WithCommandFn(cmdPronounsImportFn),
	)

	cmdPronounsList = message.NewCommand(
		message.WithCommandPermissions(discord.PermissionSendMessages),
		message.WithCommandFn(cmdPronounsListFn),
	)

	pronounsCommands = message.NewCommandMap(
		message.WithMapScope(true),
		message.WithMapCommand("add", cmdPronounsAdd),
		message.WithMapCommand("remove", cmdPronounsRemove),
		message.WithMapCommand("create", cmdPronounsCreate),
		message.WithMapCommand("delete", cmdPronounsDelete),
		message.WithMapCommand("init", cmdPronounsInit),
		message.WithMapCommand("import", cmdPronounsImport),
		message.WithMapCommand("list", cmdPronounsList),
	)

	defaultPronouns = []string{
		"he/him",
		"she/her",
		"they/them",
		"they/any",
		"they/he",
		"they/she",
	}
)

func cmdPronounsFn(ctx context.Context, s *message.Session) error {
	return pronounsCommands.ExecuteCommand(ctx, s)
}

func cmdPronounsAddFn(ctx context.Context, s *message.Session) error {
	b := getBot(s)

	prn := s.Args[0]
	dp, err := modelsx.GetPronoun(ctx, s.Tx, s.Msg.GuildID, prn)
	if err != nil {
		return err
	}

	if dp == nil {
		return s.Reply(ctx, "The requested pronoun role do not exist on this guild")
	}

	if err := b.AddRole(ctx, s.Msg.GuildID, s.Msg.Author.ID, dp.RoleSnowflake, "Automated; pronoun role requested"); err != nil {
		return err
	}

	return s.Replyf(ctx, "Thou hath been branded %s", strings.ToUpper(prn))
}

func cmdPronounsRemoveFn(ctx context.Context, s *message.Session) error {
	b := getBot(s)

	prn := s.Args[0]
	dp, err := modelsx.GetPronoun(ctx, s.Tx, s.Msg.GuildID, prn)
	if err != nil {
		return err
	}

	if dp == nil {
		return s.Reply(ctx, "The requested pronoun role does not exist on this guild")
	}

	if err := b.RemoveRole(ctx, s.Msg.GuildID, s.Msg.Author.ID, dp.RoleSnowflake, "Automated; pronoun removal requested"); err != nil {
		return err
	}

	return s.Replyf(ctx, "Thou no longer bears the brand of %s", dp.Pronoun)
}

func cmdPronounsCreateFn(ctx context.Context, s *message.Session) error {
	b := getBot(s)

	prn := s.Args[0]
	exists, err := modelsx.PronounExists(ctx, s.Tx, s.Msg.GuildID, prn)
	if err != nil {
		return err
	}

	if exists {
		return s.Reply(ctx, "The given pronouns already have a role on this server")
	}

	role, err := b.CreateRole(ctx, s.Msg.GuildID, "Automated; creation of pronoun role", discord.WithRoleName(prn))
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

	return s.Reply(ctx, "Pronoun role has been created")
}

func cmdPronounsDeleteFn(ctx context.Context, s *message.Session) error {
	b := getBot(s)

	prn := s.Args[0]
	id, deleted, err := modelsx.DeletePronoun(ctx, s.Tx, s.Msg.GuildID, prn)
	if err != nil {
		return err
	}

	if !deleted {
		return s.Reply(ctx, "The requested pronoun role does not exist on this guild")
	}

	if err := b.DeleteRole(ctx, s.Msg.GuildID, id, "Automated; pronoun removal requested"); err != nil {
		return err
	}

	return s.Reply(ctx, "Deleted pronoun role from guild")
}

func cmdPronounsListFn(ctx context.Context, s *message.Session) error {
	list, err := modelsx.Pronouns(ctx, s.Tx, s.Msg.GuildID)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return s.Reply(ctx, "I'm not managing any pronoun roles for this guild")
	}

	var sb strings.Builder

	sb.WriteString("```")

	for _, p := range list {
		sb.WriteString(p.Pronoun + "\n")
	}

	sb.WriteString("```")

	return s.Reply(ctx, sb.String())
}

func cmdPronounsInitFn(ctx context.Context, s *message.Session) error {
	b := getBot(s)

	g, err := b.FetchGuild(ctx, s.Msg.GuildID)
	if err != nil {
		return err
	}

	needed := make(map[string]string)
	for _, dp := range defaultPronouns {
		needed[strings.ToLower(dp)] = ""
	}

	var sb strings.Builder
	errors := 0
	created := 0
	imported := 0
	fmtError := func(name, kind string, err error) {
		if sb.Len() == 0 {
			sb.WriteString("Errors encountered during initialization:\n```")
		}

		sb.WriteString(name + "(" + kind + "):\t" + err.Error() + "\n")
		errors++
	}

	for _, r := range g.Roles {
		_, ok := needed[strings.ToLower(r.Name)]
		if ok {
			needed[strings.ToLower(r.Name)] = r.ID
		}
	}

	for key, val := range needed {
		key = strings.ToLower(key)

		if val == "" {
			r, err := b.CreateRole(ctx, s.Msg.GuildID, "Initializing pronoun roles...", discord.WithRoleName(key))
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
		} else {
			if !snowflake.Valid(val) {
				continue
			}

			exists, err := models.Pronouns(qm.Where("role_snowflake = ?", val), qm.Where("guild_snowflake = ?", s.Msg.GuildID)).Exists(ctx, s.Tx)
			if err != nil {
				fmtError(key, "db", err)
				continue
			}

			if exists {
				continue
			}

			prn := &models.Pronoun{
				GuildSnowflake: s.Msg.GuildID,
				Pronoun:        key,
				RoleSnowflake:  val,
			}

			if err := prn.Insert(ctx, s.Tx, boil.Infer()); err != nil {
				fmtError(key, "db", err)
				continue
			}

			imported++
		}
	}

	if errors > 0 {
		sb.WriteString("```")
		return s.Reply(ctx, sb.String())
	}

	if created == 0 {
		if imported > 0 {
			return s.Replyf(ctx, "0 pronoun roles were created, but %d existing pronoun roles were imported into the database", imported)
		}

		return s.Reply(ctx, "No change was made as every default pronoun already has an associated database entry & role")
	}

	return s.Reply(ctx, "Pronoun roles initialized successfully")
}

func cmdPronounsImportFn(ctx context.Context, s *message.Session) error {
	b := getBot(s)

	g, err := b.FetchGuild(ctx, s.Msg.GuildID)
	if err != nil {
		return err
	}

	imported := 0

	for _, a := range s.Args {
		for _, r := range g.Roles {
			if strings.EqualFold(a, r.Name) ||
				snowflake.Valid(a) && strings.EqualFold(a, r.ID) {
				prn := &models.Pronoun{
					Pronoun:        strings.ToLower(r.Name),
					GuildSnowflake: s.Msg.GuildID,
					RoleSnowflake:  r.ID,
				}

				if err := prn.Insert(ctx, s.Tx, boil.Infer()); err != nil {
					return err
				}

				imported++
			}
		}
	}

	if imported > 0 {
		if imported == 1 {
			return s.Reply(ctx, "1 pronoun role was imported")
		}

		return s.Replyf(ctx, "%d pronoun roles were imported", imported)
	}

	return s.Reply(ctx, "No pronoun roles were imported")
}
