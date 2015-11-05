package commands

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/YusukeKomatsu/honoka"
)

var (
    cleanCmd = &cobra.Command{
        Use:   "cleanup",
        Short: "Cleanup no-indexed bucket data",
        Long:  "Cleanup no-indexed bucket data",
        Run: cleanCommand,
    }
)

func cleanCommand(cmd *cobra.Command, args []string) {
    cli, err := honoka.New()
    if err != nil {
        Exit(err)
    }
    list, err := cli.Clean()
    if err != nil {
        Exit(err)
    }
    if list != nil {
        for _, result := range list {
            if result.Error != nil {
                fmt.Printf("%s (Error [%v])\n", result.Bucket, result.Error)
            } else {
                fmt.Printf("%s (Success)\n", result.Bucket)
            }
        }
    } else {
        fmt.Println("Nothing to do")
    }
}

func init() {
    RootCmd.AddCommand(cleanCmd)
}