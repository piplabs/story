// Command iliad is the main entry point for the iliad consensus client.
package main

import (
	iliadcmd "github.com/storyprotocol/iliad/client/cmd"
	libcmd "github.com/storyprotocol/iliad/lib/cmd"
)

func main() {
	libcmd.Main(iliadcmd.New())
}
