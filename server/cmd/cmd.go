package cmd

import (
	"github.com/FrisovanderVeen/chat/server/server"
	"github.com/urfave/cli"
)

var Version = "test"
var helpTemplate = `NAME:
{{.Name}} - {{.Usage}}
DESCRIPTION:
{{.Description}}
USAGE:
{{.Name}} {{if .Flags}}[flags] {{end}}command{{if .Flags}}{{end}} [arguments...]
COMMANDS:
	{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
	{{end}}{{if .Flags}}
FLAGS:
	{{range .Flags}}{{.}}
	{{end}}{{end}}
VERSION:
` + Version +
	`{{ "\n"}}`

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "address, addr",
		Value: "localhost:8080",
		Usage: "http service address",
	},
}

type Cmd struct {
	*cli.App
}

//New creates a new server
func New() *Cmd {
	app := cli.NewApp()
	app.Name = "server"
	app.Author = ""
	app.Usage = "server"
	app.Description = "A server for a chat app"
	app.Flags = globalFlags

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.Action = func(c *cli.Context) error {
		addr := c.String("address")

		err := server.Run(addr)
		if err != nil {
			return err
		}
		return nil
	}

	return &Cmd{App: app}
}
