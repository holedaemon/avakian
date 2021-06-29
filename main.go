package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.sr.ht/~sircmpwn/getopt"
	"github.com/erei/avakian/internal/bot"
	"github.com/erei/avakian/internal/pkg/zapx"
	"github.com/joho/godotenv"
	"github.com/skwair/harmony"
)

func die(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func main() {
	envFile := getopt.String("e", ".env", "path to .env file")
	useEnvFile := getopt.Bool("u", false, "load variables from file?")
	debug := getopt.Bool("d", false, "run in debug mode?")

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

	token := os.Getenv("AVAKIAN_DISCORD_TOKEN")
	prefix := os.Getenv("AVAKIAN_DISCORD_PREFIX")

	client, err := harmony.NewClient(token)
	if err != nil {
		die("error creating client:", err.Error())
		return
	}

	logger := zapx.Must(*debug)
	b, err := bot.NewBot(bot.WithClient(client), bot.WithDebug(*debug), bot.WithDefaultPrefix(prefix), bot.WithLogger(logger))
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
