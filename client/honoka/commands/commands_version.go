package commands

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/YusukeKomatsu/honoka"
)

var (
    version = "0.0.1"
    versionCmd = &cobra.Command{
        Use:   "version",
        Short: "",
        Long:  "",
        Run:   versionCommand,
    }
    versionDetail bool
)

func versionCommand(cmd *cobra.Command, args []string) {
    if versionDetail {
        fmt.Printf("Honoka: %s\nClient: %s\n", honoka.Version, version)
    } else {
        fmt.Println(version)
    }
}

func init() {
    versionCmd.Flags().BoolVarP(&versionDetail, "detail", "d", false, "show detail")
    RootCmd.AddCommand(versionCmd)
}