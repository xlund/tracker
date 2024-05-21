package app

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/xlund/tracker/internal/api"
	"github.com/xlund/tracker/internal/cmdutil"
)

func Start(ctx context.Context) int {
	var port = 4000

	if os.Getenv("PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}

	db, err := cmdutil.NewDatabasePool(ctx, 16)

	if err != nil {
		println(fmt.Sprintf("Unable to connect to database: %v", err))
		return 1
	}
	defer db.Close()

	api := api.NewApi(ctx, db)
	go func() {
		api.Router().Start(fmt.Sprintf("localhost:%d", port))
	}()

	<-ctx.Done()
	api.Router().Shutdown(ctx)

	return 0
}
