package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	_ "github.com/jackc/pgx/v4/stdlib"
	"golang.org/x/oauth2/clientcredentials"

	"git.sr.ht/~sircmpwn/getopt"
	"github.com/erei/avakian/internal/avakian"
	"github.com/erei/avakian/internal/pkg/zapx"
	"github.com/joho/godotenv"
	"github.com/skwair/harmony"
)

const (
	dbMaxAttempts = 20
	dbTimeout     = 10 * time.Second
)

func die(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func main() {
	envFile := getopt.String("e", ".env", "path to .env file")
	useEnvFile := getopt.Bool("u", false, "load variables from file?")

	if err := getopt.Parse(); err != nil {
		die("error parsing command line:", err.Error())
		return
	}

	if *useEnvFile {
		if err := godotenv.Load(*envFile); err != nil {
			die("error loading .env variables:", err.Error())
			return
		}
	}

	opts := []avakian.Option{}
	prefix := os.Getenv("AVAKIAN_DISCORD_PREFIX")
	debugEnv := os.Getenv("AVAKIAN_DEBUG")
	token := os.Getenv("AVAKIAN_DISCORD_TOKEN")
	dsn := os.Getenv("AVAKIAN_DB_DSN")
	adminIDs := os.Getenv("AVAKIAN_ADMINS")

	twitterAPIKey := os.Getenv("AVAKIAN_TWITTER_API_KEY")
	twitterAPISecret := os.Getenv("AVAKIAN_TWITTER_API_SECRET")

	debug, err := strconv.ParseBool(debugEnv)
	if err != nil {
		die("unable to parse debug variable:", err)
		return
	}

	opts = append(opts,
		avakian.WithDebug(debug),
		avakian.WithLogger(zapx.Must(debug)),
	)

	if prefix != "" {
		opts = append(opts, avakian.WithDefaultPrefix(prefix))
	}

	if token != "" {
		opts = append(opts, mainClient(token))
	}

	if dsn != "" {
		opts = append(opts, mainDB(dsn))
	}

	if twitterAPIKey != "" &&
		twitterAPISecret != "" {
		opts = append(opts, mainTwitter(twitterAPIKey, twitterAPISecret))
	}

	if adminIDs != "" {
		admins := strings.Split(adminIDs, ",")
		opts = append(opts, avakian.WithAdmins(admins))
	}

	b, err := avakian.NewBot(opts...)
	if err != nil {
		die("error creating bot:", err.Error())
		return
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	if err := b.Connect(ctx); err != nil {
		die("error connecting to discord:", err.Error())
		return
	}

	<-ctx.Done()

	b.Disconnect()
}

func mainClient(token string) avakian.Option {
	client, err := harmony.NewClient(token)
	if err != nil {
		die("error creating client:", err.Error())
		return nil
	}

	return avakian.WithClient(client)
}

func mainDB(dsn string) avakian.Option {
	var db *sql.DB
	connected := false

	for i := 0; i < dbMaxAttempts && !connected; i++ {
		conn, err := sql.Open("pgx", dsn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to connect to db, trying again in %s: %s\n", dbTimeout, err.Error())
			time.Sleep(dbTimeout)
			continue
		}

		if err := conn.Ping(); err != nil {
			fmt.Fprintf(os.Stderr, "unable to ping db, trying again in %s: %s\n", dbTimeout, err.Error())
			time.Sleep(dbTimeout)
			continue
		}

		connected = true
		db = conn
	}

	return avakian.WithDB(db)
}

func mainTwitter(key, secret string) avakian.Option {
	config := &clientcredentials.Config{
		ClientID:     key,
		ClientSecret: secret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	cli := config.Client(context.Background())

	return avakian.WithTwitter(twitter.NewClient(cli))
}
