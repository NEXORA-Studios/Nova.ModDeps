# dlremote 命令使用示例

## 基本用法

### 1. 从远程服务器下载配置文件

```bash
# 从远程服务器下载配置文件
novadm dlremote https://example.com/mods/

# 强制覆盖已存在的文件
novadm dlremote https://example.com/mods/ --force

# 备份现有文件后下载
novadm dlremote https://example.com/mods/ --backup
```

### 2. 完整的远程更新流程

```bash
# 1. 从远程下载最新的配置文件
novadm dlremote https://example.com/mods/ --backup

# 2. 同步本地文件状态
novadm load package

# 3. 安装新的或更新的 mod
novadm install
```

## 远程服务器配置

### 文件结构
远程服务器需要提供以下文件结构：

```
https://example.com/mods/
├── mod.package.json
└── mod.lock.json
```

### 示例 mod.package.json
```json
{
    "__version__": 1,
    "__platform__": "modrinth",
    "minecraftVersion": "1.21.7",
    "modLoader": ["fabric"],
    "mods": [
        {
            "id": "YL57xq9U",
            "name": "Iris 1.9.1 for Fabric 1.21.7",
            "version": "l77DAK6U",
            "requiredBy": []
        },
        {
            "id": "AANobbMI",
            "name": "Sodium 0.6.13 for Fabric 1.21.6/1.21.7",
            "version": "ND4ROcMQ",
            "requiredBy": ["l77DAK6U"]
        }
    ]
}
```

### 示例 mod.lock.json
```json
{
    "lockfileVersion": 1,
    "basedir": "mods",
    "pending": [],
    "installed": [
        {
            "id": "AANobbMI",
            "version": "ND4ROcMQ",
            "path": "mods\\sodium-fabric-0.6.13+mc1.21.6.jar",
            "uri": "https://cdn.modrinth.com/data/AANobbMI/versions/ND4ROcMQ/sodium-fabric-0.6.13%2Bmc1.21.6.jar",
            "sha512": "ee97e3df07a6f734bc8a0f77c1f1de7f47bed09cf682f048ceb12675c51b70ba727b11fcacbb7b10cc9f79b283dd71a39751312b5c70568aa3ac9471407174db"
        }
    ],
    "needRemove": []
}
```

## 使用场景

### 1. 服务器 Mod 包分发
服务器管理员可以维护一个统一的 mod 配置，玩家通过 `dlremote` 命令快速同步。

### 2. 团队协作
团队成员可以共享同一个远程配置，确保所有人使用相同的 mod 版本。

### 3. 自动化更新
可以编写脚本定期从远程获取最新的 mod 配置。

### 4. 备份和恢复
使用 `--backup` 标志可以安全地更新配置，同时保留原有配置作为备份。

## 注意事项

1. **网络连接**：确保网络连接稳定，下载过程中断可能导致文件不完整
2. **文件验证**：命令会自动验证下载文件的 JSON 格式和必要字段
3. **权限问题**：确保有足够的权限写入本地文件
4. **版本兼容性**：确保远程配置与本地 Minecraft 版本兼容 