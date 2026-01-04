#!/bin/bash

# 演示如何在编译时注入版本信息

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

# 方法1: 使用 go build 手动注入
echo "=== 方法1: 使用 go build 手动注入 ==="
echo "命令: go build -ldflags \"-X internal.version.Version=$VERSION -X internal.version.Build=$BUILD -X internal.version.Commit=$COMMIT -X internal.version.Branch=$BRANCH\" -o modix_manual ./cmd/modix"
echo

go build -ldflags "-X internal.version.Version=$VERSION -X internal.version.Build=$BUILD -X internal.version.Commit=$COMMIT -X internal.version.Branch=$BRANCH" -o modix_manual ./cmd/modix

if [ $? -eq 0 ]; then
    echo "✅ 手动编译成功!"
    echo "运行 ./modix_manual --version:"
    ./modix_manual --version
    echo
    echo "运行 ./modix_manual version:"
    ./modix_manual version
else
    echo "❌ 手动编译失败"
fi

echo

# 方法2: 使用 GoReleaser (如果安装了的话)
echo "=== 方法2: 使用 GoReleaser ==="
if command -v goreleaser &> /dev/null; then
    echo "GoReleaser 已安装，可以运行: goreleaser build --snapshot"
else
    echo "GoReleaser 未安装，跳过..."
    echo "安装 GoReleaser: curl -sSfL https://raw.githubusercontent.com/goreleaser/goreleaser/master/install.sh | sh"
fi

echo

# 方法3: 演示 ldflags 语法
echo "=== 方法3: ldflags 语法说明 ==="
echo "Go 编译时使用 -ldflags 参数可以注入变量值:"
echo "  -X importpath.name=value"
echo
echo "对于我们的情况:"
echo "  -X internal.version.Version=1.0.0"
echo "  -X internal.version.Build=da4fe6b"
echo "  -X internal.version.Commit=da4fe6b118b54306ff63a04df4ab5d6e1da76b5e"
echo "  -X internal.version.Branch=main"
echo

# 清理
rm -f modix_manual

echo "=== 演示完成 ==="