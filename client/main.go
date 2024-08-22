// Command story is the main entry point for the story consensus client.
package main

import (
	storycmd "github.com/piplabs/story/client/cmd"
	libcmd "github.com/piplabs/story/lib/cmd"
)

func main() {
	libcmd.Main(storycmd.New())
}
