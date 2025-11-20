# Modix CI/CD 配置

这个项目配置了 GitHub Actions CI/CD 来自动构建和发布跨平台的二进制文件。

## 功能特性

- **自动构建**: 当代码推送到 `main` 分支时自动触发
- **跨平台支持**: 支持 Linux、macOS (Intel & Apple Silicon) 和 Windows
- **自动发布**: 自动创建 GitHub Release 并上传二进制文件
- **手动触发**: 支持通过 GitHub Actions 界面手动触发构建

## 支持的平台

| 平台           | 目标                       | 文件格式  |
| -------------- | -------------------------- | --------- |
| Linux x86_64   | `x86_64-unknown-linux-gnu` | `.tar.gz` |
| macOS x86_64   | `x86_64-apple-darwin`      | `.tar.gz` |
| macOS ARM64    | `aarch64-apple-darwin`     | `.tar.gz` |
| Windows x86_64 | `x86_64-pc-windows-gnu`    | `.zip`    |

## 工作流说明

### 构建阶段 (build)

1. **代码检出**: 检出最新的代码
2. **工具链安装**: 安装 Rust stable 工具链和目标平台
3. **缓存设置**: 缓存 Cargo 依赖以加速构建
4. **构建**: 使用 `--release` 模式构建所有目标平台
5. **测试**: 运行测试确保代码质量
6. **打包**: 将二进制文件打包为压缩包
7. **上传**: 上传构建产物作为工件

### 发布阶段 (release)

1. **下载工件**: 下载所有平台的构建产物
2. **版本生成**: 使用日期和运行编号生成版本号
3. **创建 Release**: 创建 GitHub Release 并上传所有二进制文件

## 版本号格式

版本号格式为: `vYYYYMMDD-运行编号`
例如: `v20241120-123`

## 手动触发

你可以在 GitHub Actions 界面中点击 "Run workflow" 来手动触发构建和发布。

## 配置文件

- 主要配置: `.github/workflows/release.yml`
- Rust 配置: `Cargo.toml`
- 源代码: `src/`

## 故障排除

如果构建失败，请检查:

1. **Cargo.lock**: 确保 `Cargo.lock` 文件存在且是最新的
2. **依赖项**: 检查所有依赖项是否兼容
3. **测试**: 确保所有测试都能通过
4. **权限**: 确保 GitHub Actions 有发布权限

## 使用下载的二进制文件

### Linux/macOS

```bash
# 下载并解压
curl -L <release-url>/modix-x86_64-unknown-linux-gnu.tar.gz | tar xz

# 或者手动下载解压
tar -xzf modix-x86_64-unknown-linux-gnu.tar.gz
cd modix-x86_64-unknown-linux-gnu
chmod +x modix
./modix --help
```

### Windows

```powershell
# 下载并解压 zip 文件
Expand-Archive -Path modix-x86_64-pc-windows-gnu.zip -DestinationPath .
cd modix-x86_64-pc-windows-gnu
.\modix.exe --help
```
