package core

import (
	"fmt"
	"os"

	"github.com/NEXORA-Studios/Nova.ModDeps/cli/utils"
)

type Logger struct{}

func (l *Logger) Debug(message string) {
	fmt.Printf("%s[DBG] %s\n%s", utils.ColorGreen, message, utils.ColorReset)
}

func (l *Logger) Info(message string) {
	fmt.Printf("%s[INF] %s\n%s", utils.ColorCyan, message, utils.ColorReset)
}

func (l *Logger) Warn(message string) {
	fmt.Printf("%s[WRN] %s\n%s", utils.ColorYellow, message, utils.ColorReset)
}

func (l *Logger) Error(message string) {
	fmt.Printf("%s[ERR] %s\n%s", utils.ColorRed, message, utils.ColorReset)
}

func (l *Logger) Fatal(message string) {
	fmt.Printf("%s[FAT] %s\n%s", utils.ColorRed, message, utils.ColorReset)
	os.Exit(1)
}
