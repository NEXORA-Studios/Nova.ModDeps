package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/NEXORA-Studios/Nova.ModDeps/cli/utils"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/meta"
	"github.com/spf13/cobra"
)

var bypass bool

var RemoveCmd = &cobra.Command{
	Use:   "remove <mod_id> [--bypass]",
	Short: "移除指定 Mod 及其依赖关系",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modID := args[0]
		metaFunc := meta.MetaFunctions{}

		// 检查 requiredBy
		mod, err := metaFunc.GetModById(modID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "移除失败: %v\n", err)
			os.Exit(1)
		}

		if !bypass && len(mod.RequiredBy) > 0 {
			// 黄色警告
			fmt.Fprintf(os.Stderr, "%s警告：该 Mod 被以下 Mod 依赖：[%v]，移除可能导致依赖问题，操作已中断！\n      使用 --bypass 参数强制删除%s\n", utils.ColorYellow, strings.Join(mod.RequiredBy, ", "), utils.ColorReset)
			os.Exit(1)
		}

		err = metaFunc.RemoveMod(modID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "移除失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("已成功移除 Mod %s 及其依赖关系\n", modID)
	},
}

func init() {
	RemoveCmd.Flags().BoolVarP(&bypass, "bypass", "b", false, "强制删除")
}
