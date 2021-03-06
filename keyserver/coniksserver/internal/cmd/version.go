package cmd

import (
	"fmt"

	"github.com/coniks-sys/coniks-go/internal"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of coniksserver.",
	Long:  `Print the version number of coniksserver.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("All software has versions. This is coniksserver's:")
		fmt.Println("coniksserver v" + internal.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
