package commands

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/YusukeKomatsu/honoka"
    "github.com/davecgh/go-spew/spew"
)

var (
    outdatedCmd = &cobra.Command{
        Use:   "outdated",
        Short: "retrive no-indexed cache data list.",
        Long:  "retrive no-indexed cache data list. if use clean method, delete these.",
        Run:   outdatedCommand,
    }
)

func outdatedCommand(cmd *cobra.Command, args []string) {
    cli, err := honoka.New()
    if err != nil {
        Exit(err)
    }

    list, err := cli.Outdated()
    if err != nil {
        Exit(err)
    }
    if list == nil {
        fmt.Println("No-indexed cache data is NOTHING.")
    } else {
        spew.Dump(list);
    }
    
}

func init() {
    RootCmd.AddCommand(outdatedCmd)
}