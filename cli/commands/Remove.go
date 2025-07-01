package commands

import (
	"fmt"
	"strings"

	"github.com/NEXORA-Studios/Nova.ModDeps/core/lock"
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
			logger.Warn(fmt.Sprintf("检查依赖失败: %v", err))
		}

		if !bypass && len(mod.RequiredBy) > 0 {
			// 黄色警告
			logger.Warn(fmt.Sprintf("该 Mod 被这些 Mod 依赖：[%v]，移除可能导致依赖问题，操作已中断！\n      使用 --bypass 参数强制删除", strings.Join(mod.RequiredBy, ", ")))
			return
		}

		err = metaFunc.RemoveMod(modID)
		if err != nil {
			logger.Error(fmt.Sprintf("移除失败: %v", err))
			return
		}

		// 新增：写入 lockfile needremove
		lock.AddNeedRemove(modID, mod.Version)

		logger.Info(fmt.Sprintf("已成功移除 Mod %s", modID))
	},
}

func init() {
	RemoveCmd.Flags().BoolVarP(&bypass, "bypass", "b", false, "强制删除")
}
