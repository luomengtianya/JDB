package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Version:          "1.0.0",
	TraverseChildren: true,
}

func init() {
	RootCmd.AddCommand(JGetCmd.JGetCmd)
}
