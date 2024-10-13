package cmd

import (
	"fmt"
	
	"github.com/gkwa/itsfilmnoir/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of itsfilmnoir",
	Long:  `All software has versions. This is itsfilmnoir's`,
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo := version.GetBuildInfo()
		fmt.Println(buildInfo)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
