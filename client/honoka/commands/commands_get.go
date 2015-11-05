package commands

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/YusukeKomatsu/honoka"
)

var (
    getCmd = &cobra.Command{
        Use: "get",
        Short: "get cached data, use specified key",
        Long:  "get cached data, use specified key",
        Run: getCommand,
    }
)

func getCommand(cmd *cobra.Command, args []string) {
    if len(args) == 0 {
        Exit(fmt.Errorf("Set cache keys"))
    }
    cli, err := honoka.New()
    if err != nil {
        Exit(err)
    }
    for _, key := range args {
        val, err := cli.GetJson(key)
        if err != nil {
            fmt.Printf("%s: %v\n", key, err)
        } else {
            fmt.Printf("%s: %v\n", key, string(val))
        }
    }
}

func init() {
    RootCmd.AddCommand(getCmd)
}