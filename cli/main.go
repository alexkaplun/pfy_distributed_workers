package cli

import (
	"log"
	"os"

	"github.com/alexkaplun/pfy_distributed_workers/service"
	"github.com/pkg/errors"

	"github.com/alexkaplun/pfy_distributed_workers/storage"

	"github.com/urfave/cli"
)

const (
	DB_FILENAME = "db/sqlite.db"
	MAX_WORKERS = 4
)

var logger = log.New(os.Stdout, "cli: ", log.Ldate|log.Ltime|log.Lshortfile)

func Run(args []string) bool {

	app := cli.NewApp()
	app.Commands = cli.Commands{
		{
			Name: "run",
			Subcommands: cli.Commands{
				{
					Name: "workers",
					Action: func(_ *cli.Context) error {
						// connect to the DB
						logger.Println("initiate the database")
						db, _ := storage.New(DB_FILENAME)

						//ensure there is some data OR always create a new DB
						//if err := db.MustHaveDB(); err != nil {
						if err := db.InitDB(); err != nil {
							return errors.Wrap(err, "error checking the DB existence")
						}

						service := service.New(db, MAX_WORKERS)
						// start processing the urls
						logger.Println("start processing urls")
						service.Run()
						logger.Println("service work complete")
						return nil
					},
				},
			},
		},
	}

	if err := app.Run(args); err != nil {
		logger.Printf("app failed with error: %v\n", err)
		return false
	}
	return true
}
