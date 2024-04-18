package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/xlund/tracker/internal/api"
	"github.com/xlund/tracker/internal/cmdutil"
)

func Exectue(ctx context.Context) int {
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
	srv := api.Server(port)

	go func() {
		_ = srv.ListenAndServe()
	}()

	println("Listening on localhost:" + strconv.Itoa(port))
	<-ctx.Done()
	_ = srv.Shutdown(ctx)

	return 0
}
