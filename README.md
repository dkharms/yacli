<div align="center">

  <img width="256" height="256" src="https://user-images.githubusercontent.com/29202384/230708826-fbd86910-6611-40bb-b567-b5d71dc98ad4.jpeg">

</div>

<div align="center">

  <a href="">![GitHub](https://img.shields.io/github/license/dkharms/yacli)</a>
  <a href="">![Go Report Card](https://goreportcard.com/badge/github.com/dkharms/yacli)</a>
  <a href="">![Test Workflow](https://github.com/dkharms/yacli/actions/workflows/test.yml/badge.svg)</a>

</div>

### About

`yacli` is an open source command line tool written in `GoLang` that allows developers to easily create powerful and customizable command line interfaces (`CLIs`) for their applications.
With `yacli`, you can quickly build interactive `CLI` applications that are easy to use and provide a rich user experience.

`yacli` comes with a simple and intuitive `API` that makes it easy to define commands, flags, and arguments for your `CLI`.
You can create subcommands, and define flags with short and long names.
`yacli` also supports various types of flags, such as boolean, string, integer, float and provides built-in support for input validation and error handling.

### Features

#### Aesthetics and Simplicity

Commands do not require the use of global variables or `init()` function.
Just declaratively describe your command and let `yacli` do the rest.

#### Explicit Work With Flags and Arguments

Validation comes by default - you can take a break from arguments and flags validation.

#### Help Template

Verbose help message is already shipped.

```bash
echo [ -u | -l ] message amount
Just prints <message> in format you specified

Flags:
    -h | --help [BOOL] - Print this message
    -u | --uppercase [BOOL] - Print <message> in uppercase
    -l | --lowercase [BOOL] - Print <message> in lowercase

Arguments:
    * message [STRING] - Message to print
    * amount [INTEGER] - Print <message> `n` times
```

### How To Start

#### 1. Install `yacli`

The first step is to install yacli on your system. You can do this by running the following command:
```bash
$ go get github.com/dkharms/yacli
```

#### 2. Define your `CLI` commands

```go
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

func echo(ctx yacli.Context) error { … }
```

#### 3. Build and run your `CLI` application

Once you have defined your CLI commands, you can build and run your application using the following commands:

```bash
$ go build
$ ./echo --uppercase true "I love Computer Science" 2
$ I LOVE COMPUTER SCIENCE
$ I LOVE COMPUTER SCIENCE
```

### What's Next

Checkout [docs](https://pkg.go.dev/github.com/dkharms/yacli) or [examples](https://github.com/dkharms/yacli/tree/main/examples).
