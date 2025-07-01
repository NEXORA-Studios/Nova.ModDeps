package main

import (
	"fmt"

	"github.com/NEXORA-Studios/Nova.ModDeps/cli/commands"
	"github.com/NEXORA-Studios/Nova.ModDeps/cli/utils"
	"github.com/NEXORA-Studios/Nova.ModDeps/core"
	"github.com/spf13/cobra"
)

var logger = core.Logger{}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "novadm",
		Short: "Nova 系列产品：Mod 依赖管理工具",
	}

	rootCmd.AddCommand(utils.InfoCmd)
	rootCmd.AddCommand(commands.InitCmd)
	rootCmd.AddCommand(commands.SearchCmd)
	rootCmd.AddCommand(commands.VersionCmd)
	rootCmd.AddCommand(commands.ListCmd)
	rootCmd.AddCommand(commands.AddCmd)
	rootCmd.AddCommand(commands.RemoveCmd)
	rootCmd.AddCommand(commands.InstallCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(fmt.Sprintf("执行命令失败: %v\n", err))
	}
}
