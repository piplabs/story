package main

import (
	e2ecmd "github.com/storyprotocol/iliad/e2e/cmd"
	libcmd "github.com/storyprotocol/iliad/lib/cmd"
)

func main() {
	libcmd.Main(e2ecmd.New())
}
