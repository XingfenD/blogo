# Blogo - 基于Go的轻量级博客引擎

[![License](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/XingfenD/blogo)](https://goreportcard.com/report/github.com/XingfenD/blogo)

Blogo 是一个使用 Go 语言开发的简约博客引擎，支持 Markdown 格式文章，内置 SQLite 数据库，提供响应式前端模板。

## 功能特性

- 📝 Markdown 文章支持
- 🏷️ 分类与标签系统
- 📆 时间线归档
- 🎨 响应式主题设计
- ⚡ 极速构建与渲染
- 🔒 基于文件的简单数据存储

## 快速开始

### 前置要求
- Go 1.24+
- SQLite3

### 安装步骤
```bash
# 克隆仓库
git clone https://github.com/XingfenD/blogo.git

# 进入项目目录
cd blogo

# 安装依赖
go mod tidy

# 启动服务
go run main.go
```

## 配置说明

编辑 config.toml 文件：

```toml
[basic]
port2listen = 8080         # 监听端口
base_url = 'http://localhost:8080' # 站点地址
root_path = 'website'      # 资源根目录

[user]
name = "Your Name"         # 用户名称
avatar_url = "/img/avatar.png" # 头像路径
description = "个人博客"    # 站点描述

# 更多配置项参考 config_example.toml
```

## 项目结构

```plaintext
blogo/
├── website/             # 前端资源
│   ├── template/        # HTML模板
│   ├── static/          # 静态资源
│   └── data/            # 数据库文件
├── module/              # Go模块
│   ├── router/          # 路由处理
│   ├── sqlite/          # 数据库操作
│   └── tpl/             # 模板引擎
└── config.toml          # 配置文件
```

## 技术栈

- 后端: Go 1.24
- 数据库: SQLite3
- 模板引擎: Go html/template
- Markdown渲染: Blackfriday
- 前端: HTML5/CSS3

## 许可协议

本项目采用 [Mozilla Public License 2.0](https://opensource.org/licenses/MPL-2.0) 开源协议。

本项目中使用了字节跳动图标库提供的部分图标。

## 待办事项

- [ ] 实现后台管理页面
- [ ] 文档完善
