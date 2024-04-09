package main

import (
	"context"
	"os"
	"os/signal"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xlund/tracker/internal/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	ret := cmd.Exectue(ctx)
	os.Exit(ret)
}
