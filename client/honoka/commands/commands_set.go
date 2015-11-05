package commands

import (
    "fmt"
    "strconv"
    "github.com/spf13/cobra"
    "github.com/YusukeKomatsu/honoka"
)

var (
    setCmd = &cobra.Command{
        Use:   "set [key] [value] [expire]",
        Short: "Cache new data",
        Long:  "Cache new data if specified key is not used yet or caches (use specified key) are expired.",
        Run:   setCommand,
    }
)

func setCommand(cmd *cobra.Command, args []string) {
    if len(args) < 3 {
        Exit(fmt.Errorf("Set invalid argments"))
    }
    cli, err := honoka.New()
    if err != nil {
        Exit(err)
    }
    expire, _ := strconv.ParseInt(args[2], 10, 64)
    err = cli.Set(args[0], args[1], expire)
    if err != nil {
        Exit(err)
    }
    fmt.Println("success.")
}

func init() {
    RootCmd.AddCommand(setCmd)
}