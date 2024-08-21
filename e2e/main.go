package main

import (
	e2ecmd "github.com/piplabs/story/e2e/cmd"
	libcmd "github.com/piplabs/story/lib/cmd"
)

func main() {
	libcmd.Main(e2ecmd.New())
}
