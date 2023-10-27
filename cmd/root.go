package cmd

import (
	"os"

	"github.com/manyids2/go-tools/tui/loops"
	"github.com/manyids2/go-tools/tui/views/layout"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-tools",
	Short: "Collection of tools to visualize data.",
	Long:  `Collection of tools to visualize data.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	ui := layout.NewUI("./")
	loops.Run(ui)
}

func init() {}
