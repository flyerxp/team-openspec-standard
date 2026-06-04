# Hertz + OpenSpec 团队统一 Web 架构规范（全局通用）
## 一、规范适用范围
团队所有 Golang Hertz Web 项目，全局统一架构、目录、分层、编码、Lint 规则。
与团队 Kitex RPC 规范编码统一、架构隔离、互不污染。

## 二、整体架构分层
### 分层固定顺序（Web专属）
Controller → Service → Logic / DAL / Convert

### 各层职责
1. Controller 控制器层
- Web 路由接口入口
- 负责参数绑定、基础校验、日志埋点、错误封装、响应返回
- 无任何业务逻辑、数据查询、事务处理

2. Service 业务服务层
- 一对一对应 Web 接口
- 业务流程编排、参数校验、事务控制、数据组装
- 允许直接调用 DAL 操作数据库、缓存
- 允许调用 Logic 通用方法、Convert 模型转换

3. Logic 通用逻辑层
- 全局复用业务逻辑、工具、缓存管理、计算规则
- 无状态、可单例
- 不依赖 Service、Controller、Convert
- 仅可只读查询 DAL
- 与 Kitex RPC 项目 Logic 规范完全一致

4. DAL 数据访问层
1. 目录结构
    - 严格按照数据库实例分目录：biz/dal/gormL/{db_name}/
    - 每个库独立包含：gormL、gormL/where
2. 文件命名强制规则
    - 表模型文件：{表名}.go
    - 查询构造器文件：where/{表名}.go
    - 仓库文件：{表名}.go（与模型文件对应）
3. 代码约束
    - 模型结构体必须定义 TableName 方法
    - 查询条件必须收敛到 where 构造器，禁止在 repo 拼接 SQL
    - repo 无状态，不使用全局单例
    - 每个 repo 必须实现 GetWhere() 方法，返回对应 where 实例
    - 链式查询风格完全仿照 Ent 框架，类型安全、无裸 SQL
    - 分页统一使用 limit+1 模式，统一返回分页结构体

5. Convert 模型转换层
- DB结构体 ↔ Web请求/响应结构体 双向转换
- 仅字段映射、格式化，无业务逻辑

## 三、固定目录结构（Hertz专属）
```
├── /                    # 服务入口
├── router/              # 路由注册层（Web专属）
├── biz/                 # 业务核心层
│   ├── controller/      # Web接口控制器层
│   ├── service/         # 业务服务层
│   ├── logic/           # 公共复用逻辑 & 通用工具
│   ├── dal/             # 数据访问层
│   │   ├── gormL/      # MySQL ORM 数据访问
│   │   │   └── {db_name}/   # 按【数据库实例名】分目录（多库隔离）
│   │   │      └── where/ # 类型安全链式查询构造器
│   ├── redis/   # 当前库 Redis 缓存访问
│   ├── convert/         # 数据模型转换器
│   └── test/            # 单元测试
├── conf/                # 配置层
├── hz_gen/              # 自动生成代码（禁止手动修改）
├── layout/              # 代码生成模板
├── script/              # 部署脚本
└── test/                # 集成测试
```

## 四、合法依赖（强制）
Controller → Service
Service → Logic
Service → DAL
Service → Convert
Logic → DAL(只读)
Convert → Logic

## 五、非法依赖（红线禁止）
- Logic 禁止依赖 Service / Controller / Convert
- Controller 禁止直接操作 DAL / Logic
- Convert 禁止编写业务逻辑
- 禁止跨层乱调用

## 六、统一编码规范（全团队统一，与Kitex一致）
1. 结构体必须聚合初始化 T{}
2. 变量优先 := 类型推导，禁止冗余类型声明
3. 禁止变量遮蔽（内外重名）
4. 无限循环统一 for{}
5. JSON、带引号字符串使用反引号原生字符串
6. 禁止 main 全局状态变量，状态全部结构体私有化
7. Lint 静态可优化警告行尾 // NOLINT
8. 所有注释统一 // 单行中文注释

## 七、查询构造器规范
- 基类统一封装存储、Build 构建
- 子类仅实现业务查询条件
- 内部方法私有，外部仅链式调用
- 完全仿 Ent 体验、类型安全、零硬编码

## 八、分页统一规范
- 统一分页结构体
- limit+1 判断是否有下一页
- 统一 HasMore 字段
- 统一列表裁剪逻辑

## 九、语言规范
- 所有文档、注释、说明、需求、设计全部中文
- 禁止拼音、禁止中英文混杂、禁止缩写乱用
