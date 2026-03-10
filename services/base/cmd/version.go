package cmd

import (
	"fmt"

	"github.com/sdq-codes/usegro-api/build"
	"github.com/sdq-codes/usegro-api/version"
	"github.com/spf13/cobra"
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
		fmt.Println("Build Time:\t", build.Time)
		fmt.Println("Build User:\t", build.User)
	},
}
