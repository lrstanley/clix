package main

import (
	"github.com/lrstanley/clix/v2"
)

type Flags struct {
	Name string `name:"name" default:"world" help:"name to print"`
}

var cli = &clix.CLI[Flags]{
	Application: clix.Application{
		Name:        "simple-app",
		Description: "a simple app that prints hello world",
		Links:       clix.GithubLinks("github.com/lrstanley/clix", "master", "https://liam.sh/"),
	},
}

func main() {
	cli.ParseWithDefaults()
	cli.Logger.Info("hello", "name", cli.Flags.Name)
}
