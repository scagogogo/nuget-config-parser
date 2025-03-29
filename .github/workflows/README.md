# GitHub Actions 工作流

本项目使用 GitHub Actions 自动化测试和验证过程。以下是可用的工作流及其功能。

## 1. Go CI 工作流 (ci.yml)

这个工作流在每次推送到主分支或创建 Pull Request 时运行，执行以下任务：

- 运行所有单元测试
- 生成测试覆盖率报告并上传到 Codecov
- 使用 golangci-lint 进行代码质量检查
- 编译并运行所有示例代码

### 触发条件
- 推送到 `main` 或 `master` 分支
- 创建针对 `main` 或 `master` 分支的 Pull Request
- 手动触发 (workflow_dispatch)

### 配置选项

如果需要自定义工作流的行为，可以修改 `.github/workflows/ci.yml` 文件：

- **Go 版本**: 更改 `go-version` 参数
- **超时设置**: 示例代码运行的超时时间默认为 30 秒，可以根据需要调整
- **linter 配置**: 可以通过 `args` 参数传递额外的 golangci-lint 参数

## 2. 定时测试工作流 (scheduled-tests.yml)

这个工作流每周自动运行一次，在不同的操作系统和 Go 版本上测试代码的兼容性：

- 在 Ubuntu, macOS 和 Windows 上运行测试
- 使用多个 Go 版本进行测试 (从 1.19 到 1.23.2)
- 可选：发送测试结果通知到 Slack

### 触发条件
- 每周一 UTC 时间 6:00 自动运行
- 可以手动触发

### 配置 Slack 通知

要启用 Slack 通知，需要在 GitHub 仓库中设置 `SLACK_WEBHOOK` secret：

1. 在 Slack 中创建一个 incoming webhook
2. 在 GitHub 仓库设置中添加 secret：
   - 进入 Settings > Secrets and variables > Actions
   - 添加名为 `SLACK_WEBHOOK` 的新 Repository secret，值为 webhook URL

### 自定义测试矩阵

可以在 `.github/workflows/scheduled-tests.yml` 文件中修改 `matrix` 配置来调整测试环境：

```yaml
strategy:
  matrix:
    go-version: ['1.19', '1.20', '1.21', '1.22', '1.23.2']
    os: [ubuntu-latest, macos-latest, windows-latest]
```

## 工作流状态徽章

可以将以下徽章添加到 README.md 文件中，显示工作流的当前状态：

```markdown
[![Go CI](https://github.com/scagogogo/nuget-config-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/nuget-config-parser/actions/workflows/ci.yml)
[![Scheduled Tests](https://github.com/scagogogo/nuget-config-parser/actions/workflows/scheduled-tests.yml/badge.svg)](https://github.com/scagogogo/nuget-config-parser/actions/workflows/scheduled-tests.yml)
```

## 疑难解答

### 示例超时

如果示例代码在 CI 环境中超时，可能是因为它们在等待输入或需要特定环境设置。解决方法：

1. 检查示例是否需要用户交互，如果需要，在 CI 环境中提供默认值
2. 确保任何文件系统操作使用的是临时目录
3. 增加 `.github/workflows/ci.yml` 文件中的超时设置

### 不同操作系统上的测试失败

如果测试在特定操作系统上失败，但在其他系统上通过：

1. 检查路径分隔符（Windows 使用 `\`，而 Unix 系统使用 `/`）
2. 确保文件路径比较考虑到操作系统差异
3. 对于包含环境变量或用户目录的路径，使用 `os.ExpandEnv` 或 `os.UserHomeDir` 函数
4. 添加操作系统特定的测试跳过逻辑：`if runtime.GOOS == "windows" { t.Skip("Skipping on Windows") }` 