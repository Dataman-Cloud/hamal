package main

import (
	"os"

	"github.com/Dataman-Cloud/hamal/src/cmd/command"

	//log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	hamal := cli.NewApp()
	hamal.Name = "hamal client"
	hamal.Usage = "dataman hamal command-line"
	hamal.Version = "0.1"

	hamal.Commands = []cli.Command{
		command.NewCreateCommand(),
	}
	hamal.Run(os.Args)
}
