package main

import (
	"log"
	"os"

	"github.com/Lysoul/gocommon/postgres"
	"github.com/Lysoul/todolist/app"
	"github.com/Lysoul/todolist/db/migrations"
	"github.com/urfave/cli/v2"
)

func main() {
	cli := &cli.App{
		Name:  "TodoApp",
		Usage: "A simple todo app",
		Commands: []*cli.Command{
			app.CliCommand(),
			postgres.CliCommand(migrations.Migration),
		},
	}
	err := cli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
