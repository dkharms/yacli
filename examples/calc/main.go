package main

import (
	"fmt"
	"log"

	"github.com/dkharms/yacli"
)

var root = yacli.NewRootCommand(
	yacli.WithCommandDescription("Do some basic arithemtic operations"),
	yacli.WithArguments(
		yacli.NewArgument("x", "First operand in sum operation", yacli.Integer),
		yacli.NewArgument("y", "Second operand in sum operation", yacli.Integer),
	),
	yacli.WithMutualExclusiveFlags(
		yacli.NewFlag("sum", "s", "Do '+' arithemtic operation", yacli.Bool),
		yacli.NewFlag("diff", "d", "Do '-' arithemtic operation", yacli.Bool),
	),
	yacli.WithAction(calc),

	yacli.WithSubcommand(yacli.NewCommand(
		"sum", yacli.WithCommandDescription("Calculate sum of two integers"),
		yacli.WithArguments(
			yacli.NewArgument("x", "First operand in sum operation", yacli.Integer),
			yacli.NewArgument("y", "Second operand in sum operation", yacli.Integer),
		),
		yacli.WithAction(sum),
	)),

	yacli.WithSubcommand(yacli.NewCommand(
		"diff", yacli.WithCommandDescription("Calculate difference between two integers"),
		yacli.WithArguments(
			yacli.NewArgument("x", "First operand in diff operation", yacli.Integer),
			yacli.NewArgument("y", "Second operand in diff operation", yacli.Integer),
		),
		yacli.WithAction(diff),
	)),
)

func calc(ctx yacli.Context) error {
	var res int
	x, y := ctx.Arguments().Integer("x"), ctx.Arguments().Integer("y")

	if v, isSet := ctx.Flags().Bool("sum"); isSet && v {
		res = x + y
	}

	if v, isSet := ctx.Flags().Bool("diff"); isSet && v {
		res = x - y
	}
	fmt.Println(res)

	return nil
}

func sum(ctx yacli.Context) error {
	fmt.Println(ctx.Arguments().Integer("x") + ctx.Arguments().Integer("y"))
	return nil
}

func diff(ctx yacli.Context) error {
	fmt.Println(ctx.Arguments().Integer("x") - ctx.Arguments().Integer("y"))
	return nil
}

func main() {
	if err := root.Run(); err != nil {
		log.Fatal(err)
	}
}
