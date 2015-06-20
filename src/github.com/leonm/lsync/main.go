package main

import "os"
import "github.com/codegangsta/cli"


func main() {

  app := cli.NewApp()

  app.Name = "lsync"

  app.Commands = []cli.Command{
    {
      Name:    "server",
      Aliases: []string{"S"},
      Usage:   "run an lsync server",
      Action:  newServerCommand(),
    },
  }

  app.Run(os.Args)
}
