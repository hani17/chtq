package main

import (
	"github.com/hani17/chtq/cmd"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "chtq"
	app.Description = "Communication app"
	app.Commands = []cli.Command{
		cmd.Web,
	}
	app.Run(os.Args)
}
