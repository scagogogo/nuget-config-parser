# NuGet配置解析器示例

本目录包含一系列示例，演示如何使用NuGet配置解析器库操作NuGet配置文件。

## 示例列表

1. **01_basic_parsing** - 基本解析配置文件
   * 演示如何从文件和字符串解析NuGet配置
   * 展示如何访问配置中的各个部分（包源、活跃包源、禁用的包源、配置选项）

2. **02_search_config** - 查找配置文件
   * 演示如何查找系统中的NuGet配置文件
   * 展示如何从项目目录向上查找配置文件
   * 展示如何列出所有可用的配置文件

3. **03_create_config** - 创建新的配置文件
   * 演示如何从零创建NuGet配置
   * 展示如何添加包源、活跃包源、凭证和配置选项
   * 展示如何保存配置到文件

4. **04_modify_config** - 修改现有配置
   * 演示如何加载并修改现有的NuGet配置
   * 展示如何添加、移除和修改包源
   * 展示如何更新凭证和配置选项

5. **05_package_sources** - 包源管理
   * 详细演示包源相关操作
   * 展示如何添加各种类型的包源（HTTP, 本地文件夹等）
   * 展示如何启用/禁用/移除包源

6. **06_credentials** - 凭证管理
   * 详细演示凭证相关操作
   * 展示如何添加用户名/密码凭证
   * 展示如何更新和移除凭证

7. **07_config_options** - 配置选项管理
   * 详细演示配置选项相关操作
   * 展示如何设置全局包文件夹、代理设置等
   * 展示如何管理依赖版本和包还原设置

8. **08_serialization** - 配置序列化
   * 详细演示配置序列化和反序列化
   * 展示如何将配置对象转换为XML
   * 展示如何从文件、字符串、流解析配置

## 运行示例

每个示例都是独立的Go程序，可以单独运行。要运行示例，请确保已安装Go，并在项目根目录下执行以下命令：

```bash
go run examples/01_basic_parsing/main.go
```

所有示例都会创建临时文件和目录来展示功能，并会在程序结束时自动清理。

## 注意事项

- 这些示例使用 Go 标准库和 NuGet 配置解析器库（github.com/scagogogo/nuget-config-parser）
- 所有示例都有详细的注释，解释每个步骤的操作和目的
- 示例中展示的配置文件结构遵循 NuGet 配置文件的标准格式
- 大多数示例包括输出预览，可以不运行就了解预期结果

## 重要函数和类型参考

示例中使用的主要类型和函数有：

- `nuget.NewAPI()` - 创建主API实例
- `api.ParseFromFile()` - 从文件解析配置
- `api.ParseFromString()` - 从字符串解析配置
- `api.SaveConfig()` - 保存配置到文件
- `api.SerializeToXML()` - 序列化配置为XML字符串
- `api.AddPackageSource()` - 添加包源
- `api.RemovePackageSource()` - 移除包源
- `api.SetActivePackageSource()` - 设置活跃包源
- `api.AddCredential()` - 添加凭证
- `api.RemoveCredential()` - 移除凭证
- `api.AddConfigOption()` - 添加配置选项
- `api.RemoveConfigOption()` - 移除配置选项

详细API参考请查看项目根目录下的README.md文件。 