package utils

import (
	"fmt"

	"github.com/spf13/cobra"
)

func WelcomeMessage() {
	fmt.Printf("%s_____   __                   ______________  ___\n___  | / /________   _______ ___  __ \\__   |/  /\n__   |/ /_  __ \\_ | / /  __ `/_  / / /_  /|_/ / \n_  /|  / / /_/ /_ |/ // /_/ /_  /_/ /_  /  / /  \n/_/ |_/  \\____/_____/ \\__,_/ /_____/ /_/  /_/   \n%s\n", ColorGreen, ColorReset)
	fmt.Printf("%sNova.ModDeps (NovaDM CLI) v1.0%s\n", ColorGreen, ColorReset)
	fmt.Printf("%s轻量的 Minecraft Mod 依赖管理工具%s\n\n", ColorGreen, ColorReset)
}

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "显示 Nova.ModDeps 的版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		WelcomeMessage()
	},
}
