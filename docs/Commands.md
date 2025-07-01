# 命令列表

## 基本命令

-   `novadm init` 初始化配置文件
-   `novadm search <query>` 搜索 Mod
-   `novadm version` 查看版本信息
-   `novadm list` 查看 Mod 列表
-   `novadm add <project_id> <version_id>` 添加 Mod
-   `novadm remove <project_id> <version_id>` 移除 Mod
-   `novadm install` 安装 Mod
-   `novadm load <from>` 在 package.json 和 lock.json 之间互相补充加载
-   `novadm dlremote <url>` 从远程获取配置文件
-   `novadm cache-clean` 清理缓存

## novadm init [package|lock] [--force]

在当前工作目录 (cwd) 中创建 `mod.package.json` 和 `mod.lock.json` 文件，并初始化配置文件。

-   `[package|lock]` 为可选参数，指定要初始化的文件类型：
    -   `package`: 仅初始化 `mod.package.json`
    -   `lock`: 仅初始化 `mod.lock.json`
    -   不指定参数: 初始化两个文件
-   `--force, -f`: 强制覆盖已存在的文件

> [!NOTE]
>
> -   当初始化 `mod.package.json` 时，会自动检测 Minecraft 版本和 ModLoader 类型
> -   如果文件已存在且未使用 `--force` 标志，命令会失败并提示用户

## novadm list

查看当前 Mod 列表。

## novadm search <name> [page]

使用 Modrinth API 搜索 Mod，支持分页。

-   `<name>` 为 Mod 名称
-   `[page]` 为页码，默认为 1

## novadm version <project_id>

查看指定 Modrinth 工程的所有版本信息

-   `<project_id>` 为 Modrinth 工程 ID

## novadm add <project_id> <version_id>

添加指定 Mod 及其所有依赖（递归）

-   `<project_id>` 为 Modrinth 工程 ID
-   `<version_id>` 为 Modrinth 版本 ID

> [!NOTE]
> 当添加 Mod 时，会自动检查依赖关系，并添加依赖的 Mod
> 若 Mod 作者未指定版本，且目前没有安装任意版本的依赖，则会提示先安装依赖

## novadm remove <project_id> [--bypass]

移除指定 Mod 及其依赖关系

-   `<project_id>` 为 Modrinth 工程 ID

> [!NOTE]
> 当移除 Mod 时，会自动检查依赖关系
> 如果移除的 Mod 有依赖，则会提示先移除依赖
> 你也可以使用 `--bypass` 跳过依赖检查

## novadm install

同步 mods 文件夹与 lockfile

## novadm load <from>

在 `mod.package.json` 和 `mod.lock.json` 之间互补加载

-   `<from>` 为数据源，可选值：
    -   `package`: 从 `mod.package.json` 加载到 `mod.lock.json`
    -   `lock`: 从 `mod.lock.json` 加载到 `mod.package.json`

> [!NOTE]
>
> -   当使用 `package` 时，会将 `mod.package.json` 中的所有 mod 添加到 `mod.lock.json` 中：
>     -   如果文件已存在于本地，则添加到 `installed` 列表
>     -   如果文件不存在，则添加到 `pending` 列表
> -   当使用 `lock` 时，会将 `mod.lock.json` 中 installed 和 pending 的所有 mod 添加到 `mod.package.json` 中，并自动获取 mod 名称

## novadm dlremote <url> [--force] [--backup]

从远程获取 `mod.package.json` 和 `mod.lock.json` 的内容

-   `<url>` 为远程服务器的基础 URL，例如：`https://example.com/mods/`
-   `--force, -f`: 强制覆盖已存在的文件
-   `--backup, -b`: 在覆盖前备份现有文件

> [!NOTE]
>
> -   命令会自动在 URL 后添加 `mod.package.json` 和 `mod.lock.json` 文件名
> -   下载的文件会进行 JSON 格式验证，确保文件完整性
> -   建议在下载后使用 `novadm load package` 来同步本地文件状态
> -   使用 `--backup` 标志可以在覆盖前自动备份现有配置文件

参考 [dlremote-example.md](./dlremote-example.md) 配置一个更新源

## novadm cache-clean

清理本地请求缓存
