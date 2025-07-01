# Nova.ModDeps (NovaDM CLI)

Nova.ModDeps (NovaDM CLI) 是一个用于管理 Minecraft 模组依赖的命令行工具，支持 Modrinth 平台，帮助开发者和玩家便捷地搜索、添加、移除和管理模组依赖。

## 主要功能

-   初始化模组依赖项目
-   搜索并添加 Modrinth 模组（目前尚未支持 CurseForge）
-   列出已添加的模组依赖
-   移除指定模组
-   版本管理与元数据操作

## 安装与构建

### 通过源码构建

1. 安装 [Go 1.24+](https://go.dev/dl/)
2. 克隆本仓库：
    ```bash
    git clone https://github.com/NEXORA-Studios/Nova.ModDeps.git
    cd Nova.ModDeps
    ```
3. 构建可执行文件：
    ```bash
    go build -o bin/novadm.exe cli/main.go
    ```

### 通过 GitHub Actions 产物下载

每次主分支推送或 PR 合并后，GitHub Actions 会自动为 Windows、macOS、Linux 构建产物\
可在 [Actions 页面](https://github.com/你的用户名/Nova.ModDeps/actions) 下载对应平台的 novadm.exe。

## 基本用法

[命令列表](./docs/Commmands.md)

## 贡献

欢迎 Issue 与 PR！

---

如有问题请在 GitHub 提 Issue 反馈。
