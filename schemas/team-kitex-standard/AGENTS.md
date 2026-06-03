# OpenSpec 全局 AI 代理规则（团队 Kitex 微服务专用）
## 优先级：全局最高，所有项目强制生效

# 1. 语言强制规则（MUST）
1. 所有文档、注释、需求、设计、方案必须纯中文
2. 禁止拼音、禁止中英文混杂、禁止无意义英文缩写
3. 技术关键字保留英文：ctx、db、req、resp、gorm、kitex
4. 业务语义全部使用中文描述

# 2. 框架固定架构（MUST）
技术栈：Golang + Kitex RPC + GORM + Redis
分层固定：Handler / Service / Logic / DAL / Convert

## 分层职责严格约束
### Handler（入口层）
- 只做：参数透传、日志埋点、错误包装、响应返回
- 禁止：任何业务逻辑、数据查询、数据处理

### Service（业务实现层）
- 对应每一个 RPC 接口
- 负责业务编排、参数校验、调用 DAL/Logic/Convert
- 允许直接操作 DB、Redis（通过 DAL Repo）
- 固定结构体模板：XxxService + NewXxxService + Run

### Logic（通用逻辑层）
- 存放多 Service 复用逻辑、缓存逻辑、工具方法
- 无状态、可全局单例
- 禁止依赖 Service、Handler、Convert
- 仅可只读读取 DAL 数据

### DAL（数据访问层）
- 严格按数据库名称分目录：biz/dal/{db_name}/
- 每个库独立：gorml / gorml/where / redis
- gorml：模型、Repo、CURD
- where：仿 Ent 链式查询构造器、基类继承、私有方法封装
- Repo 无状态，每次 new 新对象，不使用单例

### Convert（模型转换层）
- 只做 DB模型 ↔ RPC模型 转换
- 禁止业务逻辑、禁止判断、禁止流程编排

# 3. 目录结构强制固定
biz/dal/{db_name}/gorml
biz/dal/{db_name}/gorml/where
biz/dal/{db_name}/redis

# 4. Golang 编码强制规范（MUST）
1. 结构体必须聚合初始化：var x T / x = T{}
2. 变量赋值优先 := 自动推导，禁止冗余类型声明
3. 禁止变量遮蔽，内部作用域禁止与外部重名
4. 无限循环统一 for{}，禁止 for 1
5. JSON/含双引号字符串，必须使用反引号原生字符串
6. main.go 禁止定义全局状态变量，所有状态封装结构体私有变量
7. Lint 可静态化警告，行尾加 // NOLINT
8. 所有注释使用 //，禁止其他注释方式

# 5. 查询构造器固定规范
- 基类 BaseWhere 统一存储条件、Build 构建
- 子类只写一行 where 条件
- 外部无法调用内部私有方法
- 全程链式调用、类型安全、无裸 SQL

# 6. 分页规范
- 统一分页结构体、自动判断 HasMore
- 统一 limit+1 分页逻辑
- 统一列表裁剪分页

# 7. 禁止行为红线
- 禁止 Logic 依赖上层业务
- 禁止 Handler 写业务
- 禁止 Convert 写逻辑
- 禁止 DAL 单例常驻
- 禁止硬编码配置
