package commands

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/YusukeKomatsu/honoka"
    "github.com/davecgh/go-spew/spew"
)

var (
    cleanCmd = &cobra.Command{
        Use:   "clean",
        Short: "cleanup unindexed bucket data",
        Long:  "cleanup unindexed bucket data",
        Run: cleanCommand,
    }
)

func cleanCommand(cmd *cobra.Command, args []string) {
    cli, err := honoka.New()
    if err != nil {
        Exit(err)
    }
    result, err := cli.Clean()
    if err != nil {
        Exit(err)
    }
    if result != nil {
        spew.Dump(result)
    }
    fmt.Println("success")
}

func init() {
    RootCmd.AddCommand(cleanCmd)
}