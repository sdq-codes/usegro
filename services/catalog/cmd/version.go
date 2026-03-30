package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usegro/services/catalog/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:\t", version.Version)
		fmt.Println("GitCommit:\t", version.GitCommit)
		fmt.Println("Build Time:\t", version.BuildTime)
		fmt.Println("Build User:\t", version.BuildUser)
	},
}
