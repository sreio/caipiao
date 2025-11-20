package main

import (
	"os"
	"path/filepath"
	"runtime"
)

// getDatabasePath 获取数据库路径
// 桌面应用将数据库存储在用户目录
func getDatabasePath() string {
	var appDir string

	// 根据操作系统获取用户配置目录
	if runtime.GOOS == "windows" {
		// Windows: C:\Users\<user>\AppData\Roaming\caipiao
		appDir = filepath.Join(os.Getenv("APPDATA"), "caipiao")
	} else if runtime.GOOS == "darwin" {
		// macOS: ~/Library/Application Support/caipiao
		homeDir, _ := os.UserHomeDir()
		appDir = filepath.Join(homeDir, "Library", "Application Support", "caipiao")
	} else {
		// Linux: ~/.config/caipiao
		homeDir, _ := os.UserHomeDir()
		appDir = filepath.Join(homeDir, ".config", "caipiao")
	}

	// 确保目录存在
	os.MkdirAll(appDir, 0755)

	return filepath.Join(appDir, "lottery.db")
}
