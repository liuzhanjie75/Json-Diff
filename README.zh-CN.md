# jsondiff

[English](README.md)

`jsondiff` 是一个使用 Go 编写的 JSON 对比命令行工具，支持递归字段级差异、
数组移动检测、对象匹配、JSON 路径过滤和终端彩色输出。

## 功能

- 递归对比对象、数组、标量和 `null`
- 使用不同颜色显示新增、删除、修改和移动
- 基于 LCS 的数组移动检测
- 可选的递归数组顺序忽略比较
- 使用 Jaccard 键集合相似度自动匹配数组对象
- 使用 `--key` 精确匹配数组对象
- 使用 `--path` 过滤 JSON 子路径
- 支持文件路径和内联 JSON 输入
- 使用 `json.Number` 保留数字精度

## 环境要求

- Go 1.26.4 或更高版本

## 构建

Windows：

```powershell
.\build.bat
```

Linux 和 macOS：

```bash
chmod +x build.sh
./build.sh
```

脚本可以从任意当前目录调用，并始终在项目根目录执行构建。`build.bat` 更新
`jsondiff.exe`；`build.sh` 默认生成 `jsondiff`，Go 目标系统为 Windows 时生成
`jsondiff.exe`。

等价的 Go 命令为：

```bash
go build -o jsondiff .
```

Windows 下使用 `go build -o jsondiff.exe .`。

## 使用

```bash
# 对比两个文件
jsondiff old.json new.json

# 对比内联 JSON
jsondiff '{"a":1}' '{"a":2}'

# 只对比指定子路径
jsondiff old.json new.json --path "database.connection"

# 使用精确标识字段匹配数组对象
jsondiff old.json new.json --key "id"

# 将数组按无序多重集合比较
jsondiff old.json new.json --ignore-array-order

# 控制终端颜色
jsondiff old.json new.json --color always
jsondiff old.json new.json --color never
```

当参数对应一个已存在的文件时，`jsondiff` 会优先读取文件，再判断是否为内联
JSON，因此 `2024.json` 这类文件名可以正常使用。输入必须只包含一个有效 JSON
值；尾随文本或额外 JSON 值会被拒绝。

## 参数

| 参数 | 说明 | 默认值 |
|---|---|---|
| `--path` | 只对比指定的 GJSON 路径 | 整个文档 |
| `--key` | 按字段值精确匹配数组对象 | 相似度匹配 |
| `--ignore-array-order` | 对比数组时忽略元素顺序 | `false` |
| `--color` | 颜色模式：`auto`、`always` 或 `never` | `auto` |

设置 `--key` 后，未匹配对象不会再使用相似度匹配。键值会保留 JSON 类型，
因此字符串 `"1"` 与数字 `1` 不会被视为相同键值。

使用 `--ignore-array-order` 时，数组按递归多重集合语义比较：

- `[1, 2]` 与 `[2, 1]` 视为相同。
- 重复元素的数量仍然有意义。
- 嵌套数组同样忽略顺序。
- 不输出移动差异，只保留修改、新增和删除。

## 退出码

| 退出码 | 含义 |
|---|---|
| `0` | 两个 JSON 值完全相同 |
| `1` | 发现差异 |
| `2` | 输入、解析、路径或参数错误 |

## 输出示例

```text
[CHANGED]  $.version  : "1.0.0"  →  "2.0.0"
[ADDED]    $.config.verbose  : true
[REMOVED]  $.config.retries  : 3
[MOVED]    $.features[2]  [0] → [2]  : "cache"
```

## 匹配行为

数组对比分为三个阶段：

1. 对规范化 JSON 哈希运行 LCS，匹配未修改元素。
2. 对剩余对象使用 `--key` 字段精确匹配；未配置键时，使用 Jaccard 键集合
   相似度匹配。
3. 将剩余且值相同但索引不同的元素标记为移动，最后将未匹配值报告为新增或
   删除。

自动相似度匹配阈值为 `0.5`。

启用 `--ignore-array-order` 后，递归无序匹配会替代 LCS 和移动检测阶段。
默认模式仍然保留数组顺序语义。

## 开发

```bash
go test ./...
go vet ./...
```

测试文件与对应实现一起放在 `internal/` 的各包目录中。`diff` 测试使用同包
测试，因此可以直接验证私有算法，不需要在生产代码中保留测试包装接口。

## 依赖

- [Cobra](https://github.com/spf13/cobra)：命令行框架
- [fatih/color](https://github.com/fatih/color)：终端颜色
- [GJSON](https://github.com/tidwall/gjson)：路径提取

## 许可证

MIT
