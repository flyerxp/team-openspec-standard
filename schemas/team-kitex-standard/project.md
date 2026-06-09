# project\.md（Kitex RPC 服务专属工程规范）

## 优先级：全局最高，所有 Kitex RPC 微服务项目强制生效

**一句话定位**：project\.md = Kitex RPC 项目代码怎么写、工程怎么建、规范怎么守

本文档为**纯 Kitex RPC 微服务专属工程规范**，完全对齐字节跳动 Kitex 官方 Protobuf 脚手架工程标准与团队统一 DAL 规范，仅约束同步 RPC 服务工程结构、脚手架生成规则、代码分层、编码规则、DAL 数据层、分页逻辑、依赖管控，为团队所有 Kitex RPC 项目唯一统一标准。

## 一、项目基础信息

### 1\.1 技术栈

Golang \+ Kitex\(Protobuf\) \+ GORM \+ Redis \+ Protobuf

### 1\.2 适用范围

所有基于 Kitex 脚手架 Protobuf 模式开发的微 RPC 服务，仅适配同步 RPC 接口服务，不包含任何异步任务场景。

### 1\.3 规范基准

完全对齐：字节跳动 Kitex Protobuf 脚手架官方规范 \+ 团队统一 DAL 工程规范

### 1\.4 脚手架生成命令（项目固定）

项目统一使用 Protobuf 模式生成代码，下方为**标准通用模板命令**，用于根据指定 proto 协议文件自动生成 RPC 服务代码，非固定初始化命令；开发中按需替换目标 proto 文件路径即可，所有接口迭代、代码更新均以此模板为准：

```Plain Text
kitex -type protobuf  -template-dir layout -I idl/ -module github.com/flyerxp/shopping.report/v2  idl/目标服务文件.proto
```

---

## 二、全局强制规约（MUST）

### 2\.1 语言规范

1\. 所有文档、注释、业务描述、设计方案、接口说明**全部使用纯中文**；

2\. 禁止拼音、中英文混杂、无意义英文缩写；

3\. 仅保留通用技术英文关键字：`ctx`、`db`、`req`、`resp`、`gorm`、`redis`、`kitex`；

4\. 所有业务语义、变量释义、注释说明必须使用中文。

### 2\.2 通用编码规范

1\. 结构体统一聚合初始化：`var x T` / `x = T{}`；

2\. 变量优先使用 `:=` 自动类型推导，禁止冗余类型声明；

3\. 禁止变量遮蔽，内外作用域不允许出现同名变量；

4\. 无限循环统一 `for{}`，禁止 `for 1`；

5\. JSON、含双引号字符串必须使用反引号原生字符串；

6\. 入口文件、启动文件禁止定义全局业务状态变量，所有状态封装为结构体私有变量；

7\. Lint 可忽略告警行尾统一添加 `// NOLINT`；

8\. 所有注释统一使用 `//` 单行注释；

9\. 接口方法、参数校验、错误返回严格遵循 Kitex Protobuf 官方接口规范；

10\. 代码目录、文件命名严格贴合 Kitex 分层规则，禁止代码散落、目录混乱。

---

## 三、Kitex 脚手架标准目录结构（Protobuf 专属）

本章节为**脚手架自动生成 \+ 团队自定义**完整目录，所有 RPC 项目目录完全统一，禁止私自改动目录结构。

```Plain Text
# Kitex Protobuf 模式完整项目目录
project-root/
├── idl/                     # 所有 protobuf 协议文件存放目录（纯手写、不生成）
├── kitex_gen/               # Kitex 自动生成代码目录（完全自动、禁止手动修改）
├── layout/                  # 脚手架自定义模板目录（固定不变）
├── conf/                    # 项目全局配置文件目录
├── biz/
│   ├── service/             # RPC接口实现层（单次接口专属逻辑、缓存、请求入口）
│   ├── logic/               # 通用可复用业务逻辑层（公共业务、复用逻辑）
│   ├── convert/             # 模型转换层（DB模型 ↔ Protobuf模型）
│   └── dal/                 # 数据访问层
│       ├── gormL/
│       │   └── {db_name}/
│       │       ├── where/   # 链式查询构造器
│       │       └── *.go     # 表模型、Repo仓储
│       └── redis/
│           └── {db_name}/   # Redis数据层
├── go.mod
├── go.sum

```

### 3\.1 核心目录权责说明

#### 3\.1\.1 idl 目录（纯手写、版本受控）

1\. 存放项目所有 `.proto` 协议定义文件；

2\. 所有服务、方法、结构体、枚举、错误码**仅在此处定义**；

3\. 禁止在生成代码内手动修改协议结构；

4\. 迭代更新接口必须先更新 proto 文件，再重新执行脚手架生成代码。

#### 3\.1\.2 kitex\_gen 目录（全自动、禁止手动修改）

1\. 由 Kitex 脚手架根据 proto 文件**全自动生成**；

2\. 包含服务定义、客户端、服务端、protobuf 结构体、序列化代码；

3\. **禁止任何手动修改**，手动修改会被下次生成覆盖；

4\. 业务代码仅可 **引用**，不可编辑。

#### 3\.1\.3 conf 目录（全局配置唯一入口）

1\. 存放项目所有配置：服务配置、数据库配置、Redis配置、MQ配置、业务开关；

2\. 目录规范：**仅允许存放 \.yaml 配置文件，禁止存放任何 \.go 源码文件**，实现配置与代码完全隔离。

#### 3\.1\.4 biz/service 目录（RPC 接口实现层）

1\. 由脚手架自动生成，存放 RPC 服务接口实现方法；

2\. 职责：接收 RPC 请求、参数校验、编排调用 logic/convert/dal 能力、组装返回值；

3\. 允许编写**当前接口专属轻量业务逻辑、参数组装、缓存读写**；

4\. 可复用、通用、重型业务逻辑禁止写在 service 层，必须下沉至 biz/logic；禁止直接手写 SQL 操作数据库；

---

## 四、分层职责严格约束（Kitex RPC 专属）

### 4\.1 分层核心职责

#### 4\.1\.1 biz/service（RPC 接口实现层）

1\. 承接 Kitex 生成的服务接口，实现所有 RPC 方法；

2\. 统一做请求参数合法性校验、基础参数兜底；

3\. 统一捕获接口异常、打印链路日志、返回标准 RPC 错误；

4\. 仅处理**当前RPC接口单次专属逻辑**，通用业务逻辑全部下沉 logic 层；

5\. 禁止编写可复用通用业务逻辑，避免代码冗余；禁止直接写 SQL 操作数据库。

#### 4\.1\.2 biz/logic（通用可复用业务层）

1\. 存放**跨接口通用、可复用、重型业务逻辑**，统一收拢公共能力；

2\. 职责：封装可复用业务流程、聚合DAL/Redis能力、公共数据处理、复杂业务编排；

3\. 无状态设计，可全局复用；

4\. 仅依赖 DAL、Redis、Convert，**禁止反向依赖 service 上层**；

5\. 所有多个接口共用的业务逻辑必须下沉至此层，禁止在 service 层重复编写。

#### 4\.1\.3 DAL 数据访问层

1\. 严格按数据库名分目录管理，多库完全隔离；

2\. gormL：存放数据库表模型、Repo 仓储、通用 CURD 能力；

3\. where：存放链式查询构造器，统一封装查询条件；

4\. redis：存放 Redis 缓存读写逻辑；

5\. Repo 无状态，每次使用必须通过构造方法新建实例，禁止全局单例。

#### 4\.1\.4 Convert 模型转换层

1\. 仅负责：**DB数据库模型 ↔ Kitex Protobuf 业务模型** 双向转换；

2\. 禁止编写任何业务判断、流程编排、数据处理逻辑；

3\. 职责单一，仅做模型字段映射、默认值补齐。

### 4\.2 分层依赖白名单（RPC 强制唯一）

严格禁止跨层、反向依赖，白名单外所有依赖全部违规

\- biz/service → logic、Convert、DAL、Redis缓存

\- biz/logic → Convert、DAL、Redis缓存（仅下层依赖，无上层依赖）

---

## 五、DAL 层完整强制规范（对齐 Kitex 标准）

### 5\.1 查询构造器（where）规范

1\. 基类`BaseWhere` 统一存储查询条件、提供 Build 构建 SQL 能力；

2\. 子类仅实现单行条件逻辑，职责单一，内部方法私有、外部不可调用；

3\. 全程链式调用、类型安全、**禁止手写裸SQL**；

4\. 完全对齐团队 Kitex RPC 项目统一标准。

#### 5\.1\.1 Where 命名规则

1\. 文件名：`表名.go`，与 gormL 目录下表模型文件一一对应；

2\. 结构体名：`大驼峰表名 + ListWhere`。

#### 5\.1\.2 Where 标准模板

```Plain Text
package where

import "github.com/flyerxp/lib/v2/middleware/gormL"

type DemoListWhere struct {
	*gormL.BaseWhere
}

func (w *DemoListWhere) TitleLike(title string) *DemoListWhere {
	w.Where("title LIKE ?", title+"%")
	return w
}

```

### 5\.2 Repo 仓储层强制规范

1\. 每个 Repo 必须实现 `GetWhere()` 方法，返回对应 Where 实例；

2\. Repo 禁止全局单例，必须提供构造方法，每次使用新建实例；

3\. 所有 DB 操作必须绑定上下文 ctx，禁止裸 DB 实例操作；

4\. 所有查询必须通过 where 链式构造器实现，禁止裸 SQL；

5\. **Where 对象唯一获取入口强制约束**：业务代码中需要使用对应Where查询对象时，**禁止直接 new/手动实例化Where结构体**，必须通过对应Repo的 `GetWhere()` 方法获取Where实例，全程统一入口，保证查询对象初始化统一、规范统一。

#### 5\.2\.1 标准 Repo \+ 分页代码示例

```Plain Text

package shop
import (
	"context"
	"github.com/flyerxp/lib/v2/middleware/gormL"
	"github.com/flyerxp/manage/v2/golang/biz/dal/gormL/shop/where"
	"github.com/flyerxp/globalStruct/widget"
	"gorm.io/gorm"
)
// DemoInfo 示例表GORM模型
type DemoInfo struct {
	Id     int    `gorm:"column:id;primaryKey;autoIncrement"`
	Title  string `gorm:"column:title;size:255"`
	Status int    `gorm:"column:status"`
	Path string `gorm:"column:path;size:255"`
	RootId int    `gorm:"column:root_id"`
}
// TableName 指定表名
func (DemoInfo) TableName() string {
	return "demo_info"
}
// DemoRepo 数据仓储层
type DemoRepo struct{}

// NewDemoRepo 构造方法：每次新建实例，禁止全局单例
func NewDemoRepo() *DemoRepo {
	return &DemoRepo{}
}

// GetWhere 强制实现：Where对象唯一获取入口
func (r *DemoRepo) GetWhere() *where.DemoListWhere {
	return &where.DemoListWhere{BaseWhere: &gormL.BaseWhere{}} 
}

// getGormModel 获取绑定上下文的DB实例
func (r *DemoRepo) getGormModel(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
        // 有事务：创建新会话避免污染外部事务
        return tx.Session(&gorm.Session{}).Model(&DemoInfo{})
    }
    // 无事务：使用默认连接
    return gormL.GetDB(ctx).Model(&DemoInfo{})
}
// DemoListPage 统一分页结构体
type DemoListPage struct {
	List []DemoInfo
	Page widget.Page
}
// DoPage 统一分页裁剪、计算是否有下一页
func (p *DemoListPage) DoPage() *DemoListPage {
	p.Page.HasMore = len(p.List) > p.Page.Size
	if p.Page.HasMore {
		p.List = p.List[:p.Page.Size]
	}
	return p
}

// ListPage 分页查询示例方法
func (r *DemoRepo) ListPage(ctx context.Context, w *where.DemoListWhere, page, limit int) (*DemoListPage, error) {
	var list []DemoInfo
	pageObj := DemoListPage{
		List: list,
		Page: widget.Page{Page: page, Size: limit},
	}
	// 基础查询
	db := r.getGormModel(ctx, nil)
	// 拼接where条件
	if w != nil {
		db = w.Build(db)
	}
	// 排序
	db = db.Order("id desc")
	// limit+1 查询，用于判断是否有下一页
	offset := (page - 1) * limit
	db = db.Offset(offset).Limit(limit + 1)
	if err := db.Find(&list).Error; err != nil {
		return pageObj.DoPage(), err
	}
	pageObj.List = list
	return pageObj.DoPage(), nil
}


func (r *DemoRepo) UpdatePathById(ctx context.Context, id int, path string, rootId int, tx *gorm.DB) error {
	db := r.getGormModel(ctx, tx)
	return db.Where("id = ?", id).Updates(map[string]interface{}{
		"path":    path,
		"root_id": rootId,
	}).Error
}
func (r *DemoRepo) UpdatePathByIdMust(ctx context.Context, id int, path string, rootId int, tx *gorm.DB) error {
	db := r.getGormModel(ctx, tx)
	result := db.Where("id = ?", id).Updates(map[string]interface{}{
		"path":    path,
		"root_id": rootId,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *DemoRepo) Save(ctx context.Context, info *DemoInfo, tx *gorm.DB) error {
	return r.getGormModel(ctx, tx).Save(info).Error
}
```

#### 5\.2\.2 更新方法选型规范

针对路径更新场景，提供两种更新方法，需根据业务场景严格选型：

1. **普通更新方法：****`UpdatePathById`**

    - 仅返回数据库操作错误，不校验更新行数

    - 适用于非核心、允许更新失败（如数据已被删除不影响业务）的普通场景

2. **强制更新方法：****`UpdatePathByIdMust`**

    - 除数据库操作错误外，额外校验更新行数（`RowsAffected`）

    - 若未找到目标记录（更新行数为 0），直接返回 `gorm.ErrRecordNotFound` 错误

    - **强制要求**：所有必须确保数据更新成功的核心场景，**必须使用此方法**

        - 典型场景：消息消费场景、交易流水场景、核心业务状态更新场景

        - 目的：避免因记录不存在导致的静默更新失败，保障核心数据一致性
      
3. **强制要求: 除事务操作场景外，禁止在 Service、Logic 层直接调用 gormL.GetDB(ctx) 操作数据库 **

   - 所有常规数据库 CRUD 操作，必须下沉至 Dal/GormL 层统一封装，业务层只允许调用 Dal/GormL 方法。

   - 禁止业务层裸写 DB 读写操作，避免会话条件污染、数据库入口混乱、不统一管控等问题。

   - 仅**需要原子性保障、失败可回滚的完整业务事务场景**，允许业务层 通过 gormL.GetDB(ctx) 获取、操作、传递 `tx` 事务实例，以此支撑多操作的统一提交/回滚。

### 5\.3 全局统一分页规范（RPC 通用）

1\. 统一分页结构体，内置 `List` 数据列表 \+ `widget.Page` 分页信息；

2\. 统一 `limit+1` 查询方式，精准判断是否存在下一页；

3\. 通过 `DoPage()` 方法自动裁剪多余数据、自动赋值 `HasMore`；

4\. RPC 所有列表、批量查询接口必须严格遵循此规范。

### 5\.4 数据库初始化规范（init\.go）

`biz/dal/gormL/init.go` 为 RPC 项目全局唯一数据库入口，禁止分散初始化连接。

#### 5\.4\.1 核心约束

1\. 支持多数据库并行初始化，统一聚合至 `DBClient` 结构体管理；

2\. 全局单次初始化，避免重复创建数据库连接；

3\. 所有 DB 操作必须绑定 ctx 上下文，禁止裸 DB 操作；

4\. 内部自动兜底初始化，业务层无需手动调用 Init 方法；

5\. 按业务库拆分独立获取方法，职责隔离。

#### 5\.4\.2 标准 init\.go 示例

```Plain Text
package gormL

import (
	"context"
	"github.com/flyerxp/lib/v2/logger"
	"github.com/flyerxp/lib/v2/middleware/gormL"
	"gorm.io/gorm"
)

type DBClient struct {
	Shop   *gorm.DB
	Report *gorm.DB
}

var dbClient *DBClient

func Init(ctx context.Context) error {
	if dbClient != nil {
		return nil
	}
	shopDB, err := gormL.GetEngine(ctx, "readshop")
	if err != nil {
		return err
	}
	reportDB, err := gormL.GetEngine(ctx, "report")
	if err != nil {
		return err
	}
	dbClient = &DBClient{Shop: shopDB, Report: reportDB}
	return nil
}

func GetShopDB(ctx context.Context) *gorm.DB {
	if dbClient == nil {
		initCtx := logger.GetContext(context.Background(), "gormInit")
		if err := Init(initCtx); err != nil {
			panic(err)
		}
	}
	return dbClient.Shop.WithContext(ctx)
}

func GetReportDB(ctx context.Context) *gorm.DB {
	if dbClient == nil {
		initCtx := logger.GetContext(context.Background(), "gormInit")
		if err := Init(initCtx); err != nil {
			panic(err)
		}
	}
	return dbClient.Report.WithContext(ctx)
}
```

---

## 六、Kitex 脚手架生成规范（强制）

### 6\.1 生成原则

1\. 所有 RPC 服务代码、接口定义、客户端代码**必须通过官方脚手架基于 proto 文件生成**，项目初始化、新增接口、迭代更新接口均使用该方式；

2\. 禁止手写服务注册、接口方法、proto 结构体；

3\. 每次更新 proto 文件后，必须重新执行生成命令，同步更新代码；

4\. 固定使用 `-type protobuf` 模式，统一 Protobuf 协议体系。

### 6\.2 脚手架固定参数释义

\- `-type protobuf`：指定使用 Protobuf 协议生成服务代码；

\- `-template-dir layout`：使用项目统一自定义模板，保证全项目目录风格一致；

\- `-I idl/`：指定 proto 依赖检索目录；

\- `-module`：对齐项目 go\.mod 模块路径；

\- 最后参数为**当前需要生成代码的目标服务proto文件路径**，按需动态替换，不固定。

### 6\.3 生成后代码约束

1\. **kitex\_gen**：纯自动生成，只读不改；

2\. **biz/service**：自动生成空方法模板，业务仅在此实现逻辑；

3\. 新增接口、修改字段、删除方法，必须先改 proto，再重新生成。

---

## 七、RPC 项目强制红线禁止行为（MUST NOT）

1\. 禁止 Logic 层依赖上层 service 接口层，禁止反向跨层依赖；

2\. 禁止 Convert 层编写任何业务判断、数据处理、流程编排逻辑；

3\. 禁止 DAL Repo 全局单例常驻，必须每次新建实例；

4\. 禁止业务代码硬编码配置、时间、状态值、ID 等常量，统一放 conf；

5\. 禁止手写裸 SQL，所有查询必须使用 where 链式构造器；

6\. 禁止跨层、反向非法依赖，严格遵守分层白名单；

7\. 禁止变量遮蔽、冗余类型声明、非标准循环与字符串写法；

8\. 禁止入口文件定义全局业务状态变量；

9\. 禁止手动修改 kitex\_gen 自动生成目录代码；

10\. 禁止脱离 idl/proto 定义，私自新增 RPC 接口与字段。

> （注：文档部分内容可能由 AI 生成）
