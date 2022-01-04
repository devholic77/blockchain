package cmd

import (
	"fmt"

	"github.com/devholic77/duckcoin/explorer"
	"github.com/devholic77/duckcoin/rest"
	"github.com/spf13/cobra"
)

var port int = 5000
var restCmd = &cobra.Command{
	Use:   "rest [-p port]",
	Short: "start rest api server",
	Run: func(cmd *cobra.Command, args []string) {
		rest.Start(port)
	},
}
var explorerCmd = &cobra.Command{
	Use:   "exp [-p port]",
	Short: "start explorer server",
	Run: func(cmd *cobra.Command, args []string) {
		explorer.Start(port)
	},
}
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of duckcoin",
	Long:  `All software has versions. This is duckcoin's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version v0.1")
	},
}

func Execute() {
	var rootCmd = &cobra.Command{Use: "duckcoin"}
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 5000, "server port")
	rootCmd.AddCommand(versionCmd, restCmd, explorerCmd)
	rootCmd.Execute()
}
