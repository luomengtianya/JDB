package main

import (
	"jdb/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
