# 安装指南

本指南将引导您完成 NuGet Config Parser 库的安装过程。

## 系统要求

### Go 版本

- **最低要求**: Go 1.19 或更高版本
- **推荐**: Go 1.21 或更高版本

您可以通过运行以下命令检查您的 Go 版本：

```bash
go version
```

如果您需要安装或升级 Go，请访问 [官方 Go 网站](https://golang.org/dl/)。

### 操作系统支持

该库支持所有 Go 支持的主要操作系统：

- **Linux** (所有主要发行版)
- **macOS** (10.15 或更高版本)
- **Windows** (Windows 10 或更高版本)
- **FreeBSD**
- **其他 Unix 系统**

## 安装方法

### 方法 1: 使用 Go Modules（推荐）

这是安装库的推荐方法。Go Modules 会自动处理依赖管理。

#### 在现有项目中安装

如果您已经有一个 Go 项目：

```bash
# 导航到您的项目目录
cd your-project

# 添加库作为依赖
go get github.com/scagogogo/nuget-config-parser
```

#### 创建新项目

如果您正在开始一个新项目：

```bash
# 创建新目录
mkdir my-nuget-project
cd my-nuget-project

# 初始化 Go 模块
go mod init my-nuget-project

# 添加库
go get github.com/scagogogo/nuget-config-parser
```

### 方法 2: 直接克隆（开发用）

如果您想为库做贡献或需要最新的开发版本：

```bash
# 克隆仓库
git clone https://github.com/scagogogo/nuget-config-parser.git

# 导航到目录
cd nuget-config-parser

# 安装依赖
go mod download
```

## 验证安装

创建一个简单的测试文件来验证安装：

```go
// main.go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    // 创建 API 实例
    api := nuget.NewAPI()
    
    // 创建默认配置
    config := api.CreateDefaultConfig()
    
    fmt.Printf("成功！创建了包含 %d 个包源的配置\n", 
        len(config.PackageSources.Add))
    
    // 显示默认包源
    for _, source := range config.PackageSources.Add {
        fmt.Printf("包源: %s -> %s\n", source.Key, source.Value)
    }
}
```

运行测试：

```bash
go run main.go
```

如果您看到类似以下的输出，说明安装成功：

```
成功！创建了包含 1 个包源的配置
包源: nuget.org -> https://api.nuget.org/v3/index.json
```

## 依赖管理

### 查看依赖

查看项目的所有依赖：

```bash
go list -m all
```

### 更新依赖

更新到最新版本：

```bash
go get -u github.com/scagogogo/nuget-config-parser
```

更新到特定版本：

```bash
go get github.com/scagogogo/nuget-config-parser@v1.2.3
```

### 清理依赖

移除未使用的依赖：

```bash
go mod tidy
```

## 版本管理

### 使用特定版本

您可以安装特定版本的库：

```bash
# 安装特定版本
go get github.com/scagogogo/nuget-config-parser@v1.0.0

# 安装最新的预发布版本
go get github.com/scagogogo/nuget-config-parser@latest

# 安装特定分支
go get github.com/scagogogo/nuget-config-parser@main
```

### 版本约束

在您的 `go.mod` 文件中，您可以指定版本约束：

```go
module my-nuget-project

go 1.19

require (
    github.com/scagogogo/nuget-config-parser v1.0.0
)
```

## 构建配置

### 构建标签

该库支持不同的构建标签以启用或禁用特定功能：

```bash
# 使用调试信息构建
go build -tags debug

# 使用生产优化构建
go build -tags production -ldflags "-s -w"
```

### 交叉编译

为不同平台构建：

```bash
# 为 Linux 构建
GOOS=linux GOARCH=amd64 go build

# 为 Windows 构建
GOOS=windows GOARCH=amd64 go build

# 为 macOS 构建
GOOS=darwin GOARCH=amd64 go build
```

## 故障排除

### 常见问题

#### 问题 1: "module not found" 错误

```
go: github.com/scagogogo/nuget-config-parser@latest: module github.com/scagogogo/nuget-config-parser: not found
```

**解决方案:**
1. 检查您的网络连接
2. 验证仓库 URL 是否正确
3. 确保您的 Go 版本支持模块

#### 问题 2: Go 版本不兼容

```
go: github.com/scagogogo/nuget-config-parser requires go >= 1.19
```

**解决方案:**
升级您的 Go 安装到 1.19 或更高版本。

#### 问题 3: 代理问题

如果您在公司防火墙后面：

```bash
# 设置 Go 代理
export GOPROXY=https://proxy.golang.org,direct

# 或禁用代理
export GOPROXY=direct

# 设置私有模块
export GOPRIVATE=github.com/your-company/*
```

#### 问题 4: 权限问题

在某些系统上，您可能需要调整权限：

```bash
# 确保 Go 路径可写
chmod -R 755 $GOPATH

# 或使用 sudo（不推荐）
sudo go get github.com/scagogogo/nuget-config-parser
```

### 获取帮助

如果您遇到安装问题：

1. **检查文档**: 查看 [在线文档](https://scagogogo.github.io/nuget-config-parser/)
2. **搜索问题**: 在 [GitHub Issues](https://github.com/scagogogo/nuget-config-parser/issues) 中搜索
3. **创建问题**: 如果找不到解决方案，创建新的 issue
4. **社区支持**: 查看 GitHub Discussions

## 开发环境设置

如果您计划为库做贡献：

### 1. Fork 和克隆

```bash
# Fork 仓库（在 GitHub 上）
# 然后克隆您的 fork
git clone https://github.com/YOUR-USERNAME/nuget-config-parser.git
cd nuget-config-parser
```

### 2. 设置开发依赖

```bash
# 安装开发工具
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装测试依赖
go mod download
```

### 3. 运行测试

```bash
# 运行所有测试
go test ./...

# 运行带覆盖率的测试
go test -cover ./...

# 运行基准测试
go test -bench=. ./...
```

### 4. 代码质量检查

```bash
# 运行 linter
golangci-lint run

# 格式化代码
goimports -w .

# 检查模块
go mod verify
go mod tidy
```

## 生产部署

### 构建优化

为生产环境构建优化的二进制文件：

```bash
# 构建优化的二进制文件
go build -ldflags "-s -w" -o myapp

# 使用 UPX 进一步压缩（可选）
upx --best myapp
```

### Docker 部署

创建 Dockerfile：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w" -o myapp

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/myapp .
CMD ["./myapp"]
```

构建和运行：

```bash
docker build -t my-nuget-app .
docker run my-nuget-app
```

## 下一步

安装完成后：

1. 阅读 [快速开始指南](./quick-start.md)
2. 查看 [API 参考](/zh/api/)
3. 探索 [示例](/zh/examples/)
4. 了解 [最佳实践](./best-practices.md)

## 更新和维护

### 定期更新

定期检查更新：

```bash
# 检查可用更新
go list -u -m all

# 更新到最新版本
go get -u github.com/scagogogo/nuget-config-parser

# 清理依赖
go mod tidy
```

### 安全更新

关注安全公告并及时更新：

```bash
# 检查已知漏洞
go list -json -deps | nancy sleuth

# 或使用 govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

通过遵循本安装指南，您应该能够成功安装和设置 NuGet Config Parser 库。如果您遇到任何问题，请参考故障排除部分或寻求社区帮助。
