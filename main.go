package main

import (
	"github.com/ahdekkers/go-zipdir/zipdir"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	var pathToDir string
	var outPath string
	rootCmd := &cobra.Command{
		Use:   "zipdir",
		Short: "zipdir path/to/dir",
		Long:  "zipdir is a tool to zip a directory",
		Run: func(cmd *cobra.Command, args []string) {
			err := zipdir.ZipToDir(pathToDir, outPath)
			if err != nil {
				cmd.Printf("Failed to zip directory: %v", err)
			}
		},
	}

	rootCmd.Flags().StringVarP(&pathToDir, "in", "i", "",
		"Path to dir which should be zipped")
	rootCmd.Flags().StringVarP(&outPath, "out", "o", "",
		"Output zip file including path")

	err := rootCmd.MarkFlagRequired("in")
	if err != nil {
		log.Printf("[WARN]  Failed to mark input flag as required flag: %v\n", err)
	}
	err = rootCmd.MarkFlagDirname("in")
	if err != nil {
		log.Printf("[WARN]  Failed to mark input flag as directory: %v\n", err)
	}
	err = rootCmd.MarkFlagRequired("out")
	if err != nil {
		log.Printf("[WARN]  Failed to mark output flag as required flag: %v\n", err)
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Printf("[ERROR] Failed to run zipdir: %v\n", err)
	}
}
