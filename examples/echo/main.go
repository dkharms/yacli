package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/dkharms/yacli"
)

var root = yacli.NewRootCommand(
	yacli.WithCommandDescription("Just prints <message> in format you specified"),
	yacli.WithMutualExclusiveFlags(
		yacli.NewFlag("uppercase", "u", "Print <message> in uppercase", yacli.Bool),
		yacli.NewFlag("lowercase", "l", "Print <message> in lowercase", yacli.Bool),
	),
	yacli.WithArguments(
		yacli.NewArgument("message", "Message to print", yacli.String),
		yacli.NewArgument("amount", "Print <message> `n` times", yacli.Integer),
	),
	yacli.WithAction(echo),
)

func echo(ctx yacli.Context) error {
	n := ctx.Arguments().Integer("amount")
	message := ctx.Arguments().String("message")

	if v, isSet := ctx.Flags().Bool("uppercase"); isSet && v {
		message = strings.ToUpper(message)
	}

	if v, isSet := ctx.Flags().Bool("lowercase"); isSet && v {
		message = strings.ToLower(message)
	}

	for i := 0; i < n; i++ {
		fmt.Println(message)
	}

	return nil
}

func main() {
	if err := root.Run(); err != nil {
		log.Fatal(err)
	}
}
