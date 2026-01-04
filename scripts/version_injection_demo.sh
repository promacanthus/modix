#!/bin/bash

# GoReleaser 版本信息注入完整演示

echo "=== GoReleaser 版本信息注入演示 ==="
echo

# 获取当前信息
VERSION="1.0.0"
BUILD=$(git rev-parse --short HEAD 2>/dev/null || echo "dev-build")
COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "dev-commit")
BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "main")

echo "当前构建信息:"
echo "  Version: $VERSION"
echo "  Build: $BUILD"
echo "  Commit: $COMMIT"
echo "  Branch: $BRANCH"
echo

# 关键：使用引号包围路径
echo "=== 正确的 ldflags 语法 ==="
echo "注意：路径必须用引号包围"
echo

LD_FLAGS="-X 'github.com/promacanthus/modix/internal/version.Version=$VERSION' \
-X 'github.com/promacanthus/modix/internal/version.Build=$BUILD' \
-X 'github.com/promacanthus/modix/internal/version.Commit=$COMMIT' \
-X 'github.com/promacanthus/modix/internal/version.Branch=$BRANCH'"

echo "编译命令:"
echo "go build -ldflags \"$LD_FLAGS\" -o modix_injected ./cmd/modix"
echo

# 编译
go build -ldflags "$LD_FLAGS" -o modix_injected ./cmd/modix

if [ $? -eq 0 ]; then
    echo "✅ 编译成功!"
    echo
    echo "运行 ./modix_injected --version:"
    ./modix_injected --version
    echo
    echo "运行 ./modix_injected version:"
    ./modix_injected version
else
    echo "❌ 编译失败"
fi

echo

# 清理
rm -f modix_injected test_version test_version.go

echo "=== 关键要点 ==="
echo "1. 使用完整的模块路径: github.com/promacanthus/modix/internal/version"
echo "2. 路径必须用引号包围: -X 'path.name=value'"
echo "3. GoReleaser 会自动处理引号和变量替换"
echo "4. ldflags 语法: -X importpath.variable=value"
echo

echo "GoReleaser 配置示例:"
echo "ldflags:"
echo "  - -s -w"
echo "  - -X 'github.com/promacanthus/modix/internal/version.Version={{.Version}}'"
echo "  - -X 'github.com/promacanthus/modix/internal/version.Build={{.ShortCommit}}'"
echo "  - -X 'github.com/promacanthus/modix/internal/version.Commit={{.Commit}}'"
echo "  - -X 'github.com/promacanthus/modix/internal/version.Branch={{.Branch}}'"
echo

echo "=== 演示完成 ==="