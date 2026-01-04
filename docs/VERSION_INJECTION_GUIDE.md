# GoReleaser 版本信息注入指南

## 概述

在 Go 项目中，我们经常需要在编译时注入版本信息，如版本号、构建时间、Git 提交哈希等。GoReleaser 提供了强大的功能来自动化这个过程。

## 关键概念

### 1. ldflags 语法

Go 编译器的 `-ldflags` 参数允许在链接阶段设置变量值：

```bash
go build -ldflags "-X importpath.variable=value" -o binary ./cmd/app
```

### 2. 变量要求

要被 `-X` 标志设置的变量必须：
- 是 `string` 类型
- 在包级别声明
- 可以被赋值（不是常量）

## 实现步骤

### 1. 创建版本信息包

创建 `internal/version/version.go`：

```go
package version

// Version information injected by GoReleaser
var (
	Version = "dev"
	Build   = "dev"
	Commit  = "none"
	Branch  = "main"
)
```

### 2. 在主程序中使用

```go
package main

import (
	"fmt"
	"github.com/your-module/internal/version"
)

func main() {
	fmt.Printf("Version: %s\n", version.Version)
	fmt.Printf("Build: %s\n", version.Build)
	fmt.Printf("Commit: %s\n", version.Commit)
	fmt.Printf("Branch: %s\n", version.Branch)
}
```

### 3. GoReleaser 配置

在 `.goreleaser.yaml` 中配置 `ldflags`：

```yaml
builds:
  - id: app-linux-amd64
    main: ./cmd/app
    binary: app
    env:
      - CGO_ENABLED=0
    goos: linux
    goarch: amd64
    ldflags:
      - -s -w
      - -X 'github.com/your-module/internal/version.Version={{.Version}}'
      - -X 'github.com/your-module/internal/version.Build={{.ShortCommit}}'
      - -X 'github.com/your-module/internal/version.Commit={{.Commit}}'
      - -X 'github.com/your-module/internal/version.Branch={{.Branch}}'
```

## 关键要点

### 1. 路径格式

- **必须使用完整的模块路径**：`github.com/your-module/internal/version`
- **路径必须用引号包围**：`-X 'path.name=value'`
- **原因**：路径中的斜杠会被 shell 解析，需要用引号保护

### 2. GoReleaser 变量

GoReleaser 提供的内置变量：
- `{{.Version}}` - 版本号
- `{{.ShortCommit}}` - 短提交哈希（7位）
- `{{.Commit}}` - 完整提交哈希
- `{{.Branch}}` - 当前分支名
- `{{.Tag}}` - 当前标签
- `{{.Date}}` - 构建日期

### 3. 手动编译示例

```bash
# 获取构建信息
VERSION="1.0.0"
BUILD=$(git rev-parse --short HEAD)
COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)

# 编译并注入版本信息
go build -ldflags "
  -X 'github.com/your-module/internal/version.Version=$VERSION'
  -X 'github.com/your-module/internal/version.Build=$BUILD'
  -X 'github.com/your-module/internal/version.Commit=$COMMIT'
  -X 'github.com/your-module/internal/version.Branch=$BRANCH'
" -o app ./cmd/app
```

## 常见问题

### 1. 版本信息没有被注入

**原因**：
- 路径不正确
- 路径没有用引号包围
- 变量不是 `string` 类型

**解决**：
- 使用 `go list -f '{{.ImportPath}}' ./internal/version` 检查正确路径
- 确保路径用引号包围
- 确保变量是 `string` 类型

### 2. 编译错误

**原因**：
- ldflags 语法错误
- 路径中的特殊字符

**解决**：
- 使用单引号包围整个 `-X` 参数
- 检查路径是否正确

## 最佳实践

1. **使用专门的版本包**：避免污染主包
2. **提供默认值**：开发时使用合理的默认值
3. **使用引号**：始终用引号包围路径
4. **测试注入**：确保版本信息正确显示
5. **文档化**：记录版本信息的格式和用途

## 示例输出

```bash
$ ./app --version
app version 1.0.0
Build: da4fe6b
Commit: da4fe6b118b54306ff63a04df4ab5d6e1da76b5e
Branch: main

$ ./app version
App Version Information

Version: 1.0.0
Build: da4fe6b
Commit: da4fe6b118b54306ff63a04df4ab5d6e1da76b5e
Branch: main

Environment Information
Go Version: go1.25.5
Platform: darwin/arm64
Compiler: gc
```

## 参考

- [GoReleaser 官方文档](https://goreleaser.com/)
- [Go 编译器 ldflags 文档](https://pkg.go.dev/cmd/link)
- [Go 模块路径最佳实践](https://go.dev/doc/modules/)