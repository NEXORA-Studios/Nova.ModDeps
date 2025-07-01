package main

import (
	"fmt"
	"os"

	"github.com/NEXORA-Studios/Nova.ModDeps/cli/commands"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "novadm",
		Short: "Nova 系列产品：Mod 依赖管理工具",
	}

	rootCmd.AddCommand(commands.InitCmd)
	rootCmd.AddCommand(commands.SearchCmd)
	rootCmd.AddCommand(commands.VersionCmd)
	rootCmd.AddCommand(commands.ListCmd)
	rootCmd.AddCommand(commands.AddCmd)
	rootCmd.AddCommand(commands.RemoveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
