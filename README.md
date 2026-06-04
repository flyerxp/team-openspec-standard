# 团队 Go 微服务 OpenSpec 全局规范仓库 README

## 📌 仓库简介

本仓库为**团队内部统一 Golang 微服务 OpenSpec 全局标准仓库**，是所有 Kitex RPC 微服务、Hertz Web 服务的代码规范、架构标准、AI 生成规则、工程目录的唯一权威来源。

彻底解决多项目规范混乱、框架规则互相污染、AI 生成代码风格不统一、新旧项目架构不一致、团队编码风格割裂的问题，实现**一套通用编码规范、两套框架隔离架构**的标准化落地。

## ✨ 核心价值

- **终结重复规范配置**：无需每个项目单独编写架构、目录、编码、AI 规则，全局统一继承

- **双框架完美隔离**：Kitex RPC / Hertz Web 架构分层完全独立，互不污染

- **全团队风格统一**：通用编码、注释、Lint、数据层规范全局一致

- **AI 生成标准化**：强制 AI 区分框架生成对应代码，杜绝乱分层、乱建目录

- **一键初始化模板**：内置 project\.md 标准模板，新项目 init 直接对齐团队规范

- **版本可管控同步**：规范统一版本管理，全员拉取即可同步最新标准

## 📁 仓库目录结构

```plain
schemas/
├── team-kitex-standard/        # Kitex RPC 微服务专属全局规范
│   ├── config.yaml             # 框架匹配、分层、编码强制校验规则
│   ├── AGENTS.md               # AI 代码生成最高优先级指令
│   ├── spec.md                 # 完整架构+目录+编码规范文档
│   ├── project.md              # 项目级默认规约（初始化模板）
│   └── README.md               # Kitex 规范接入教程
│
├── team-hertz-standard/        # Hertz Web 服务专属全局规范
│   ├── config.yaml             # 框架匹配、分层、编码强制校验规则
│   ├── AGENTS.md               # AI 代码生成最高优先级指令
│   ├── spec.md                 # 完整架构+目录+编码规范文档
│   ├── project.md              # 项目级默认规约（初始化模板）
│   └── README.md               # Hertz 规范接入教程
│
└── README.md                   # 仓库总说明文档（本文档）
```

## 🔍 双框架隔离机制（核心特性）

两套规范自动识别项目框架、差异化生效，**通用规则统一、架构规则隔离**。依托 OpenSpec 内置 `project.md` 项目模板能力，新项目初始化自动套用团队统一规约，无需人工重复配置。

### 1\. 全局通用（所有项目强制遵守）

- 全部文档、注释、业务描述强制中文，禁止拼音、中英文混杂

- 结构体统一聚合初始化，杜绝空结构体告警

- 变量使用自动类型推导，禁止冗余类型声明

- 禁止变量遮蔽、统一无限循环写法 `for{}`

- JSON 字符串统一使用原生反引号，禁止转义符拼接

- main 文件禁止全局状态变量，状态统一结构体封装

- Lint 告警统一 `// NOLINT` 规范抑制

- DAL 层统一分库架构：`dal/gormL/{db_name}/where`

- 统一链式查询构造器、分页逻辑、数据层编码规范

### 2\. 框架专属隔离规则

#### ✅ Kitex RPC 专属（Hertz 不生效）

- 分层架构：**Handler → Service → Logic/DAL/Convert**

- 适配 IDL、kitex\_gen RPC 体系

- RPC 服务专属结构体模板、业务编排规则

#### ✅ Hertz Web 专属（Kitex 不生效）

- 分层架构：**Controller → Service → Logic/DAL/Convert**

- 适配 router 路由、hz\_gen Web 体系

- Web 接口参数绑定、路由注册专属规范

## 🏗️ 统一工程架构标准

所有项目强制遵循**五层分层架构**，职责单一、依赖可控：

1. **入口层**（Handler/Controller）：仅参数透传、日志、错误包装，无业务逻辑

2. **业务层 Service**：业务流程编排、事务控制、数据源调用

3. **通用层 Logic**：全局复用逻辑、工具方法、缓存规则，无状态不依赖上层

4. **数据层 DAL**：分库管理 MySQL/Redis，统一 CURD、链式查询

5. **转换层 Convert**：仅做模型映射，无任何业务逻辑

## ⚙️ 部署与接入指南

### 1\. 全局部署（仅需执行一次）

将仓库内两套规范文件夹，放置到 OpenSpec 全局公共目录：

```bash
~/.local/share/openspec/schemas/
```
```cmd
%USERPROFILE%\.openspec\schemas\
```


执行校验命令，确认规范生效：

```bash
openspec validate
```

### 2\. 项目接入方式

在项目 `openspec/config.yaml` 中添加继承配置：

- **Kitex RPC 项目**

- `extends: team-kitex-standard`

- **Hertz Web 项目**

- `extends: team-hertz-standard`

接入后自动继承对应框架所有架构、编码、AI 生成规范。

## ✅ 规范生效范围

- OpenSpec 需求、设计、文档强制中文规范

- AI 全自动代码生成风格统一、架构不混乱

- 代码分层、目录结构、依赖关系强制校验

- Golang 编码风格、Lint 告警统一规范

- 团队代码评审、项目初始化唯一标准

## 📝 维护说明

- 所有全局规范迭代、规则更新统一提交至本仓库

- 团队成员拉取最新代码后，重新放入全局目录即可完成全员同步

- 禁止单项目私自定制通用规范，统一从全局仓库更新

## 🖥️ 技术栈适配

**适配栈**：Golang \+ Kitex/Hertz \+ GORM \+ Redis \+ OpenSpec

**适用项目**：团队所有 Go 微服务、Web 服务、RPC 服务

> （注：文档部分内容可能由 AI 生成）
