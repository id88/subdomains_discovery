# Build and Release Workflow 配置说明

这个 GitHub Actions 工作流会自动编译和打包你的 Go 项目。

## 功能特性

### 1. 自动构建 (`build.yml`)

- **多平台支持**: 自动为以下平台编译:
  - Linux (amd64, arm64)
  - Windows (amd64)
  - macOS (amd64, arm64)

- **触发条件**:
  - 推送到 `main` 或 `master` 分支
  - 创建以 `v` 开头的标签 (例如: `v1.0.0`)
  - Pull Request 到 `main` 或 `master` 分支
  - 手动触发 (workflow_dispatch)

- **构建优化**:
  - 使用 `-trimpath` 移除文件系统路径
  - 使用 `-ldflags="-s -w"` 减小二进制文件大小
  - 禁用 CGO (`CGO_ENABLED=0`) 以获得静态二进制文件
  - 自动注入版本号和构建时间

- **打包格式**:
  - Windows: `.zip` 文件
  - Linux/macOS: `.tar.gz` 文件
  - 每个包都包含: 可执行文件、字典文件、README 和 LICENSE

### 2. 自动发布

当你推送一个标签时 (例如 `v1.0.0`), 工作流会:
1. 构建所有平台的二进制文件
2. 生成 SHA256 校验和
3. 自动创建 GitHub Release
4. 上传所有编译好的包到 Release

### 3. 代码质量检查 (`lint.yml`)

- **代码检查**: 使用 golangci-lint 进行代码质量检查
- **安全扫描**: 使用 gosec 进行安全漏洞扫描

### 4. 测试

- 自动运行所有测试
- 生成代码覆盖率报告
- 上传到 Codecov (可选)

## 使用方法

### 发布新版本

```bash
# 1. 提交所有更改
git add .
git commit -m "准备发布 v1.0.0"

# 2. 创建标签
git tag -a v1.0.0 -m "Release version 1.0.0"

# 3. 推送标签
git push origin v1.0.0
```

GitHub Actions 会自动:
- 编译所有平台的二进制文件
- 创建 Release
- 上传所有文件

### 手动触发构建

1. 进入 GitHub 仓库
2. 点击 "Actions" 标签
3. 选择 "Build and Release" 工作流
4. 点击 "Run workflow"

## 构建产物

构建完成后，你会得到以下文件:

```
subdomains_discovery-linux-amd64.tar.gz
subdomains_discovery-linux-arm64.tar.gz
subdomains_discovery-windows-amd64.zip
subdomains_discovery-darwin-amd64.tar.gz
subdomains_discovery-darwin-arm64.tar.gz
checksums.txt
```

## 自定义配置

### 修改 Go 版本

在 `.github/workflows/build.yml` 中修改:

```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.24.2'  # 修改这里
```

### 添加更多平台

在 `matrix` 部分添加:

```yaml
strategy:
  matrix:
    goos: [linux, windows, darwin, freebsd]  # 添加更多操作系统
    goarch: [amd64, arm64, 386]              # 添加更多架构
```

### 修改编译参数

在 `Build` 步骤中修改 `go build` 命令:

```bash
go build -v -trimpath \
  -ldflags="-s -w -X 'main.Version=${{ github.ref_name }}'" \
  -o "${OUTPUT_NAME}" .
```

## 注意事项

1. **首次使用**: 确保你的仓库有 `main` 或 `master` 分支
2. **权限**: GitHub Actions 需要有写入权限来创建 Release
3. **标签格式**: Release 只在标签以 `v` 开头时创建 (例如: `v1.0.0`, `v2.1.3`)
4. **构建时间**: 多平台构建可能需要几分钟时间

## 故障排除

### 构建失败

检查:
- Go 版本是否正确
- 依赖是否都能正常下载
- 代码是否有编译错误

### Release 未创建

确保:
- 标签以 `v` 开头
- 推送了标签到远程仓库
- GitHub Actions 有足够的权限

### 下载失败

检查:
- 网络连接
- Go 模块代理设置
- 依赖是否可访问
