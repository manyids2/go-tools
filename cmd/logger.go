package cmd

import (
	"fmt"

	"github.com/manyids2/go-tools/tui/models/logger"
	"github.com/spf13/cobra"
)

var datadir, fileext string

// loggerCmd represents the logger command
var loggerCmd = &cobra.Command{
	Use:   "logger",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		m := logger.New(datadir, fileext)
		go m.SetLogFiles()
		<-m.Loaded
		fmt.Println(m)
	},
}

func init() {
	rootCmd.AddCommand(loggerCmd)

	loggerCmd.PersistentFlags().StringVarP(&datadir,
		"datadir", "d", "./log",
		"Path to log directory")

	loggerCmd.PersistentFlags().StringVarP(&fileext,
		"fileext", "e", ".log",
		"Extension of log files")
}
