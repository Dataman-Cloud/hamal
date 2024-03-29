package main

import (
	"os"

	"github.com/Dataman-Cloud/hamal/src/hamalcli/command"
	"github.com/Dataman-Cloud/hamal/src/hamalcli/config"

	"github.com/urfave/cli"
)

func main() {
	config.InitConfig("./.hamal_cfg")

	hamal := cli.NewApp()
	hamal.Name = "hamal client"
	hamal.Usage = "dataman hamal client"
	hamal.Version = "0.1"

	hamal.Commands = []cli.Command{
		command.NewDeployCommand(),
	}
	hamal.Run(os.Args)
}
