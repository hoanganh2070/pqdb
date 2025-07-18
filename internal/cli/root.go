package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hoanganh2070/pqdb/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "PQDB",
	Short: "A Key Value Database",
	Long:  "A simple key-value database built in Go, inspired by LevelDB. It supports basic operations like put, get, and delete.",
}

func Execute() {
	c := color.New(color.FgHiMagenta, color.Bold)

	c.Print(config.Logo)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
