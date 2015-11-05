package commands

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var (
    RootCmd = &cobra.Command{
        Use:   "honoka",
        Short: "client tool for honoka (Golang file cache library)",
        Long:  "client tool for honoka (Golang file cache library)",
        Run:   func(cmd *cobra.Command, args []string) {
            versionCommand(cmd, args)
            cmd.Usage()
        },
    }
)

func Exit(err error, codes ...int) {
    var code int
    if len(codes) > 0 {
        code = codes[0]
    } else {
        code = 2
    }
    fmt.Println(err)
    os.Exit(code)
}

func Run() {
    RootCmd.Execute()
}
