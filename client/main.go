// Command iliad is the main entry point for the iliad consensus client.
package main

import (
	iliadcmd "github.com/piplabs/story/client/cmd"
	libcmd "github.com/piplabs/story/lib/cmd"
)

func main() {
	libcmd.Main(iliadcmd.New())
}
