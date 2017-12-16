package cmd

import (
	"github.com/FrisovanderVeen/chat/client/client"
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
	cli.StringFlag{
		Name:  "path",
		Value: "/echo",
		Usage: "The path to the server",
	},
}

type Cmd struct {
	*cli.App
}

//New creates a new client
func New() *Cmd {
	app := cli.NewApp()
	app.Name = "client"
	app.Author = ""
	app.Usage = "client"
	app.Description = "A client for a chat app"
	app.Flags = globalFlags

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.Action = func(c *cli.Context) error {
		addr := c.String("address")
		path := c.String("path")

		client.Run(addr, path)
		return nil
	}

	return &Cmd{App: app}
}
