# AGENTS.md — AI Agent 开发指南

> 本文件为 AI 编程助手提供项目上下文。当你在新的仓库中接手此项目时，请先阅读本文件。

## 项目概述

**jsondiff** 是一个 Go 语言编写的 CLI 工具，用于对比两个 JSON 值的差异。核心能力：
- 递归对比 JSON 对象/数组/标量/null 四种类型
- 终端彩色输出（Added=绿、Removed=红、Changed=黄、Moved=青）
- 数组 LCS 移动检测
- `--ignore-array-order` 递归无序数组比较（多重集合语义）
- 数组对象匹配：`--key` 字段值精确匹配，或未指定 key 时使用 Jaccard 相似度自动匹配
- JSON Path 过滤（`--path`）
- 支持文件路径和内联 JSON 字符串两种输入
- 默认英文文档为 `README.md`，中文文档为 `README.zh-CN.md`

## 技术栈

| 项 | 值 |
|---|---|
| 语言 | Go 1.26+ |
| 模块路径 | `github.com/zhanjie/jsondiff` |
| CLI 框架 | `github.com/spf13/cobra` |
| 终端颜色 | `github.com/fatih/color` |
| JSON Path | `github.com/tidwall/gjson` |
| JSON 解码 | 标准库 `encoding/json`，必须使用 `UseNumber()` |

## 目录结构与职责

```
main.go                              ← 入口，仅调用 cmd.Execute()
build.bat                            ← Windows 构建入口，更新根目录 jsondiff.exe
build.sh                             ← Linux/macOS 构建入口，生成平台对应二进制
install-skill.bat                    ← Windows Skill 安装器，编译后安装到 Codex skills
install-skill.sh                     ← Linux/macOS Skill 安装器
skill/jsondiff/
  SKILL.md                           ← Codex Skill 工作流
  agents/openai.yaml                 ← Skill UI 元数据
cmd/
  root.go                            ← cobra 命令定义，串联 input→jsonpath→diff→render
internal/
  input/
    input.go                         ← Resolve(arg): 自动识别文件路径 vs 内联 JSON
    input_test.go
  diff/                              ← 核心 diff 引擎（最重要的包）
    model.go                         ← 纯数据定义：Op 枚举、DiffItem 结构体
    api.go                           ← 公共 API：Options、Compare()、CompareWithOpts()
    diff.go                          ← Comparator 接口、Context、Dispatch() 调度器
    object_comparator.go             ← ObjectComparator: JSON 对象键集合运算
    array.go                         ← ArrayComparator: 数组对比三阶段编排
    match_strategy.go                ← matchByKey、matchBySimilarity、objectSimilarity
    lcs.go                           ← LCS（最长公共子序列）DP 算法
    hash.go                          ← CanonicalJSON（键排序序列化）+ HashJSON（SHA256）
    null_comparator.go               ← NullComparator: nil 值处理
    primitive_comparator.go          ← PrimitiveComparator: string/number/bool 比较
    *_test.go                        ← 与实现同包的单元测试，可直接测试私有算法
  render/
    terminal.go                      ← 彩色终端渲染 + formatValue()
    render_test.go
  jsonpath/
    extract.go                       ← 基于 gjson 的 JSON Path 子路径提取
    extract_test.go
testdata/
  simple_old.json / simple_new.json  ← 基础测试数据
```

## 核心架构

### Comparator 接口与调度

```go
// diff.go — 所有 Comparator 实现的接口
type Comparator interface {
    Compare(old, new interface{}, path string, ctx *Context) []DiffItem
}

// Context 在递归中传递配置和调度器
type Context struct {
    Opts       Options
    Dispatcher func(old, new interface{}, path string, ctx *Context) []DiffItem
}
```

**Dispatch 调度逻辑**（`diff.go`）：
1. `old==nil || new==nil` → `NullComparator`
2. 都是 `map[string]interface{}` → `ObjectComparator`
3. 都是 `[]interface{}` → `ArrayComparator`
4. 其他 → `PrimitiveComparator`

各 Comparator 需要递归时，调用 `ctx.Dispatcher()` 而非直接调用其他 Comparator，实现类型自动分派。

### 数组对比三阶段（`array.go`）

1. **Phase 1 — LCS 精确匹配**：对每个元素计算 CanonicalJSON 哈希，运行 LCS DP 算法找最长公共子序列
2. **Phase 2 — 对象匹配**：对未匹配的 object 元素选择且只选择一种策略
   - 若 `--key` 指定：按带 JSON 类型的字段值精确匹配（`matchByKey`），不再回退到相似度匹配
   - 否则：按 Jaccard 相似度匹配，阈值 0.5（`matchBySimilarity`）
3. **Phase 3 — 移动检测**：剩余未匹配但哈希相同的元素标记为 `OpMoved`
4. 最终剩余 → `OpRemoved` / `OpAdded`

### 数据流

```
cmd/root.go
  → input.Resolve(arg)           // string → interface{}
  → jsonpath.Extract(val, path)  // 可选：子路径过滤
  → diff.CompareWithOpts(...)    // interface{} → []DiffItem
  → render.Render(diffs, stdout) // []DiffItem → 终端输出
  → os.Exit(0/1)                 // 退出码
```

## 关键设计决策

### 数字精度
所有 JSON 解码必须使用 `decoder.UseNumber()`，数字类型为 `json.Number`（字符串），比较时用 `.String()` 而非 `float64`。

### 输入解析
`input.Resolve()` 必须先检查参数是否对应已存在的文件，再判断是否为内联 JSON。这保证 `2024.json`、`-old.json` 等文件名可用。JSON 解码后必须继续读取并确认 EOF，禁止接受尾随文本或多个 JSON 值。

### `--key` 精确匹配
配置 `KeyField` 时，数组对象只能按该字段匹配，不得对未匹配对象回退到 Jaccard 相似度。匹配键使用 JSON 编码保留类型，因此字符串 `"1"`、数字 `1` 和布尔值 `true` 是不同键值。

### 忽略数组顺序
`Options.IgnoreArrayOrder` 启用后，所有数组（包括嵌套数组）按多重集合比较。元素顺序不产生差异，但重复数量仍然有意义。该模式不运行 LCS 和移动检测，不得生成 `OpMoved`；剩余差异只报告修改、新增和删除。默认模式必须继续保持顺序敏感。

### CanonicalJSON（`hash.go`）
对象键按字母序排列后序列化，确保 `{"a":1,"b":2}` 和 `{"b":2,"a":1}` 哈希相同。这是数组元素精确匹配的基础。

### 相似度阈值（`array.go`）
`similarityThreshold = 0.5`：两个对象的 Jaccard 键集相似度 ≥ 0.5 时，视为"同一元素被修改"而非"删除+新增"。

### 单一职责原则
每个文件只做一件事：
- `model.go` 只有数据类型，无逻辑
- `api.go` 只有公共 API 入口
- `diff.go` 只有接口定义和调度
- 每个 Comparator 独立文件
- `lcs.go` 纯算法
- `match_strategy.go` 纯匹配策略
- `formatValue()` 在 render 包，不在 diff 包（格式化是渲染的职责）

## Coding Conventions

- **Code comments MUST be in English.** No Chinese characters in any `.go` file (source or test).
- Follow Go standard formatting (`gofmt`).
- Use table-driven tests with `t.Run()` for subtests.
- Each file should have a single responsibility — keep modules orthogonal.
- `README.md` MUST remain the default English documentation; maintain the Chinese translation in `README.zh-CN.md`.
- Keep tests beside the package they cover. Use `package diff` for private diff algorithms instead of adding production test wrappers.
- Skill installers must build successfully before replacing an existing installation. Do not commit compiled binaries under `skill/jsondiff/bin/`.

## 构建与测试

```bash
# Windows 构建
.\build.bat

# Linux/macOS 构建
./build.sh

# 运行所有测试
go test ./...

# 带详细输出
go test ./... -v

# 覆盖率
go test ./internal/diff/ -cover
```

## Windows 环境注意事项

- Go 安装在 `C:\Program Files\Go\bin\`，可能需要手动设置 PATH：
  ```powershell
  $env:Path = "C:\Program Files\Go\bin;" + $env:Path
  ```
- PowerShell 不支持 `&&`，用 `;` 分隔命令
- PowerShell 中单引号包裹的内联 JSON 传给外部程序时引号会被剥离，建议在 CMD/Git Bash 中测试内联 JSON
- 编译产物为 `jsondiff.exe`
- `build.bat` 会先从 PATH 查找 Go，再回退到 `C:\Program Files\Go\bin\go.exe`

## 退出码约定

| 码 | 含义 |
|----|------|
| 0 | JSON 完全相同 |
| 1 | 存在差异 |
| 2 | 错误 |

## 已完成的 TODO / 后续可扩展方向

- [ ] 添加更多输出格式（如 JSON 输出 `--output json`）
- [ ] 大文件流式解析（当前全部加载到内存）
- [ ] HTML diff 报告生成
- [ ] `--no-move-detect` flag 跳过大数组的移动检测
- [ ] 集成测试（编译后 exec.Command 调用）
- [ ] 为 README 补充更多高级使用示例
- [ ] CI/CD 配置（GitHub Actions）
