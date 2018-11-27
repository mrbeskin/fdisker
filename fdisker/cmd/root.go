package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mrbeskin/fdisker"
)

var rootCmd = &cobra.Command{
	Use:   "fdisker",
	Short: "fdisker is used to run interactive fdisk commands that are defined in a file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := fdisker.RunFdiskCommandFile(filePath, mountPath, writeFlag)
		if err != nil {
			fmt.Println("fdisker failed: %v", err)
		}
	},
}

var filePath string
var mountPath string
var writeFlag bool

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&writeFlag, "write", false, "set so that fdisk quits without saving")
	rootCmd.Flags().StringVarP(&filePath, "file-path", "f", "", "the path to the fdisk file commands you want to run")
	rootCmd.Flags().StringVarP(&mountPath, "mount-path", "m", "", "the path to the disk you want to run fdisk on")
	rootCmd.MarkFlagRequired("file-path")
	rootCmd.MarkFlagRequired("mount-path")
}
