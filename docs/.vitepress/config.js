import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'NuGet Config Parser',
  description: 'A Go library for parsing and manipulating NuGet configuration files',
  
  // Base URL for GitHub Pages
  base: '/nuget-config-parser/',
  
  // Language configuration
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      title: 'NuGet Config Parser',
      description: 'A Go library for parsing and manipulating NuGet configuration files',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Guide', link: '/guide/getting-started' },
          { text: 'API Reference', link: '/api/' },
          { text: 'Examples', link: '/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/nuget-config-parser' }
        ],
        sidebar: {
          '/guide/': [
            {
              text: 'Guide',
              items: [
                { text: 'Getting Started', link: '/guide/getting-started' },
                { text: 'Installation', link: '/guide/installation' },
                { text: 'Quick Start', link: '/guide/quick-start' },
                { text: 'Configuration', link: '/guide/configuration' },
                { text: 'Position-Aware Editing', link: '/guide/position-aware-editing' }
              ]
            }
          ],
          '/api/': [
            {
              text: 'API Reference',
              items: [
                { text: 'Overview', link: '/api/' },
                { text: 'Core API', link: '/api/core' },
                { text: 'Parser', link: '/api/parser' },
                { text: 'Editor', link: '/api/editor' },
                { text: 'Finder', link: '/api/finder' },
                { text: 'Manager', link: '/api/manager' },
                { text: 'Types', link: '/api/types' },
                { text: 'Utils', link: '/api/utils' },
                { text: 'Errors', link: '/api/errors' },
                { text: 'Constants', link: '/api/constants' }
              ]
            }
          ],
          '/examples/': [
            {
              text: 'Examples',
              items: [
                { text: 'Overview', link: '/examples/' },
                { text: 'Basic Parsing', link: '/examples/basic-parsing' },
                { text: 'Finding Configs', link: '/examples/finding-configs' },
                { text: 'Creating Configs', link: '/examples/creating-configs' },
                { text: 'Modifying Configs', link: '/examples/modifying-configs' },
                { text: 'Package Sources', link: '/examples/package-sources' },
                { text: 'Credentials', link: '/examples/credentials' },
                { text: 'Config Options', link: '/examples/config-options' },
                { text: 'Serialization', link: '/examples/serialization' },
                { text: 'Position-Aware Editing', link: '/examples/position-aware-editing' }
              ]
            }
          ]
        },
        socialLinks: [
          { icon: 'github', link: 'https://github.com/scagogogo/nuget-config-parser' }
        ],
        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © 2024 NuGet Config Parser'
        }
      }
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'NuGet 配置解析器',
      description: '用于解析和操作 NuGet 配置文件的 Go 库',
      themeConfig: {
        nav: [
          { text: '首页', link: '/zh/' },
          { text: '指南', link: '/zh/guide/getting-started' },
          { text: 'API 参考', link: '/zh/api/' },
          { text: '示例', link: '/zh/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/nuget-config-parser' }
        ],
        sidebar: {
          '/zh/guide/': [
            {
              text: '指南',
              items: [
                { text: '开始使用', link: '/zh/guide/getting-started' },
                { text: '安装', link: '/zh/guide/installation' },
                { text: '快速开始', link: '/zh/guide/quick-start' },
                { text: '配置', link: '/zh/guide/configuration' },
                { text: '位置感知编辑', link: '/zh/guide/position-aware-editing' }
              ]
            }
          ],
          '/zh/api/': [
            {
              text: 'API 参考',
              items: [
                { text: '概述', link: '/zh/api/' },
                { text: '核心 API', link: '/zh/api/core' },
                { text: '解析器', link: '/zh/api/parser' },
                { text: '编辑器', link: '/zh/api/editor' },
                { text: '查找器', link: '/zh/api/finder' },
                { text: '管理器', link: '/zh/api/manager' },
                { text: '类型', link: '/zh/api/types' },
                { text: '工具', link: '/zh/api/utils' },
                { text: '错误', link: '/zh/api/errors' },
                { text: '常量', link: '/zh/api/constants' }
              ]
            }
          ],
          '/zh/examples/': [
            {
              text: '示例',
              items: [
                { text: '概述', link: '/zh/examples/' },
                { text: '基本解析', link: '/zh/examples/basic-parsing' },
                { text: '查找配置', link: '/zh/examples/finding-configs' },
                { text: '创建配置', link: '/zh/examples/creating-configs' },
                { text: '修改配置', link: '/zh/examples/modifying-configs' },
                { text: '包源管理', link: '/zh/examples/package-sources' },
                { text: '凭证管理', link: '/zh/examples/credentials' },
                { text: '配置选项', link: '/zh/examples/config-options' },
                { text: '序列化', link: '/zh/examples/serialization' },
                { text: '位置感知编辑', link: '/zh/examples/position-aware-editing' }
              ]
            }
          ]
        },
        socialLinks: [
          { icon: 'github', link: 'https://github.com/scagogogo/nuget-config-parser' }
        ],
        footer: {
          message: '基于 MIT 许可证发布。',
          copyright: 'Copyright © 2024 NuGet Config Parser'
        }
      }
    }
  },
  
  themeConfig: {
    search: {
      provider: 'local'
    }
  },

  // Ignore dead links for now
  ignoreDeadLinks: true
})
