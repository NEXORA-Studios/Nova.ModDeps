# JSON 输出格式示例

## list 命令 JSON 输出

```bash
novadm list --json
```

```json
{
    "minecraftVersion": "1.21.7",
    "modLoader": ["fabric"],
    "totalMods": 2,
    "mods": [
        {
            "id": "YL57xq9U",
            "version": "l77DAK6U",
            "requiredBy": [],
            "name": "Iris",
            "versionName": "1.9.1",
            "status": "installed"
        },
        {
            "id": "AANobbMI",
            "version": "ND4ROcMQ",
            "requiredBy": ["l77DAK6U"],
            "name": "Sodium",
            "versionName": "0.6.13",
            "status": "installed"
        }
    ]
}
```

## search 命令 JSON 输出

```bash
novadm search sodium --json
```

```json
{
    "totalHits": 15,
    "currentPage": 1,
    "totalPages": 2,
    "results": [
        {
            "id": "AANobbMI",
            "name": "Sodium",
            "type": ["客户端"],
            "versionRange": ["1.16", "1.21"],
            "loaderSupport": ["fabric"],
            "description": "Sodium is a free and open-source rendering engine replacement for the Minecraft client that greatly improves frame rates and reduces micro-stutter.",
            "author": "jellysquid3",
            "categories": ["fabric", "optimization", "rendering"]
        }
    ]
}
```

## version 命令 JSON 输出

```bash
novadm version AANobbMI --json
```

```json
{
    "projectId": "AANobbMI",
    "total": 25,
    "currentPage": 1,
    "totalPages": 3,
    "versions": [
        {
            "id": "ND4ROcMQ",
            "name": "0.6.13",
            "datePublished": "2024-01-15T10:30:00Z",
            "gameVersions": ["1.21.6", "1.21.7"],
            "loaders": ["fabric"],
            "compatibility": {
                "gameMatch": true,
                "loaderMatch": true
            }
        },
        {
            "id": "ABC123XY",
            "name": "0.6.12",
            "datePublished": "2024-01-10T15:45:00Z",
            "gameVersions": ["1.21.5"],
            "loaders": ["fabric"],
            "compatibility": {
                "gameMatch": false,
                "loaderMatch": true
            }
        }
    ]
}
```

## JSON 输出字段说明

### list 命令

-   `minecraftVersion`: Minecraft 版本
-   `modLoader`: Mod 加载器类型数组
-   `totalMods`: Mod 总数
-   `mods`: Mod 列表
    -   `id`: Mod 项目 ID
    -   `version`: 版本 ID
    -   `requiredBy`: 被哪些 Mod 依赖
    -   `name`: Mod 名称
    -   `versionName`: 版本名称
    -   `status`: 状态（installed/pending/needRemove）

### search 命令

-   `totalHits`: 搜索结果总数
-   `currentPage`: 当前页码
-   `totalPages`: 总页数
-   `results`: 搜索结果列表
    -   `id`: 项目 ID
    -   `name`: 项目名称
    -   `type`: 类型（客户端/服务端）
    -   `versionRange`: 支持的版本范围
    -   `loaderSupport`: 支持的加载器
    -   `description`: 项目描述
    -   `author`: 作者
    -   `categories`: 分类标签

### version 命令

-   `projectId`: 项目 ID
-   `total`: 版本总数
-   `currentPage`: 当前页码
-   `totalPages`: 总页数
-   `versions`: 版本列表
    -   `id`: 版本 ID
    -   `name`: 版本名称
    -   `datePublished`: 发布时间
    -   `gameVersions`: 支持的游戏版本
    -   `loaders`: 支持的加载器
    -   `compatibility`: 兼容性信息
        -   `gameMatch`: 是否匹配当前游戏版本
        -   `loaderMatch`: 是否匹配当前加载器

## 使用场景

### 1. 脚本自动化

```bash
# 获取已安装的 mod 列表
novadm list --json | jq '.mods[].name'

# 搜索特定 mod
novadm search "sodium" --json | jq '.results[0].id'
```

### 2. API 集成

```bash
# 检查 mod 兼容性
novadm version AANobbMI --json | jq '.versions[] | select(.compatibility.gameMatch and .compatibility.loaderMatch)'
```

### 3. 数据导出

```bash
# 导出 mod 列表到文件
novadm list --json > mods.json

# 导出搜索结果
novadm search "fabric" --json > search_results.json
```

### 4. 监控和报告

```bash
# 统计已安装的 mod 数量
novadm list --json | jq '.totalMods'

# 检查待安装的 mod
novadm list --json | jq '.mods[] | select(.status == "pending")'
```
