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
      Usage:   "run an lsync server: lsync S <source-directory>",
      Action:  newServerCommand(),
    },
    {
      Name:    "copy",
      Aliases: []string{"C"},
      Usage:   "copy files from an lsync server: lsync C <host> <target-directory>",
      Action:  newCopyCommand(),
    },
  }

  app.Run(os.Args)
}
