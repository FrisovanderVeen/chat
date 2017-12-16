package main

import (
	"os"

	"github.com/FrisovanderVeen/chat/server/cmd"
)

func main() {
	app := cmd.New()
	app.Run(os.Args)
}
