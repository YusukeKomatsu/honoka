package commands

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/YusukeKomatsu/honoka"
)

var (
    deleteCmd = &cobra.Command{
        Use:   "delete",
        Short: "Delete cache",
        Long:  "Delete cache",
        Run:   deleteCommand,
    }
)

func deleteCommand(cmd *cobra.Command, args []string) {
    if len(args) == 0 {
        Exit(fmt.Errorf("Set cache keys"))
    }
    cli, err := honoka.New()
    if err != nil {
        Exit(err)
    }
    err = cli.Delete(args[0])
    if err != nil {
        Exit(err)
    }
    fmt.Println("success.")
}

func init() {
    RootCmd.AddCommand(deleteCmd)
}