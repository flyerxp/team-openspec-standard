# 更新后的Hertz项目规范

## 优先级：全局最高，所有 Hertz Web 前端网关 / 后台管理服务强制生效

**一句话定位**：\[project\.md\]\(project\.md\) = Hertz Web 项目代码怎么写、工程怎么建、规范怎么守
本文档为**纯 Hertz Web 服务专属工程规范**，完全对齐字节跳动 Hertz 官方 Protobuf 脚手架工程标准 \+ 团队统一 DAL 工程规范。统一 Web 项目目录结构、脚手架生成规则、分层职责、编码规范、DAL 数据层、分页逻辑、依赖管控、页面渲染规范，与 Kitex RPC 服务 DAL/Logic 规范保持一致，为团队 Hertz Web 项目唯一统一标准。

## 一、项目基础信息

### 1\.1 技术栈

Golang \+ Hertz$Protobuf$ \+ GORM \+ Redis \+ Protobuf \+ Liquid 页面渲染

### 1\.2 适用范围

所有基于 Hertz 脚手架 Protobuf 模式开发的 Web 服务、后台管理网关服务，包含 HTTP 接口服务、HTML 页面渲染服务。

### 1\.3 规范基准

完全对齐：字节跳动 Hertz 官方 Protobuf 脚手架规范 \+ 团队统一 DAL 工程规范 \+ Kitex RPC 同源分层规范

### 1\.4 脚手架生成命令（通用模板）

项目统一使用 Hertz Protobuf 模式生成代码，下方为通用模板命令，按需替换目标 proto 文件，用于更新 / 生成 HTTP 路由、handle、model 代码：

```Plain
hz update --mod=github.com/flyerxp/manage/v2 --proto_path=idl --idl=idl/目标文件.proto  --customize_package=layout/package.yaml
```

说明：`idl/shop/shop.proto` 为具体业务 Proto 协议路径，每次迭代新增接口、更新字段均替换对应 proto 文件重新生成。

---

## 二、全局强制规约（MUST）

### 2\.1 语言规范

1\. 所有文档、注释、业务描述、设计方案、接口说明**全部使用纯中文**；
2\. 禁止拼音、中英文混杂、无意义英文缩写；
3\. 仅保留通用技术英文关键字：`ctx`、`db`、`req`、`resp`、`gorm`、`redis`、`hertz`；
4\. 所有业务语义、变量释义、注释说明必须使用中文。

### 2\.2 通用编码规范

1\. 结构体统一聚合初始化：`var x T` / `x = T{}`；
2\. 变量优先使用 `:=` 自动类型推导，禁止冗余类型声明；
3\. 禁止变量遮蔽，内外作用域不允许出现同名变量；
4\. 无限循环统一 `for{}`，禁止 `for 1`；
5\. JSON、含双引号字符串必须使用反引号原生字符串；
6\. 入口文件、启动文件禁止定义全局业务状态变量；
7\. Lint 可忽略告警行尾统一添加 `// NOLINT`；
8\. 所有注释统一使用 `//` 单行注释；
9\. HTTP 接口参数校验、异常返回严格遵循 Hertz 官方规范。

---

## 三、Hertz 标准目录结构（项目强制固定）

目录完全贴合 Hertz 脚手架生成规则，区分**自动生成目录**与**手动业务目录**，禁止私自改动目录结构。

```Plain
project-root/
├── idl/                     # 所有 protobuf 协议文件（纯手写、不自动生成）
├── layout/                  # Hertz 自定义模板配置目录（固定不变）
├── conf/                    # 全局配置目录（仅存放 .yaml 配置，禁止Go文件）
├── golang/                  # Golang业务代码根目录
│   └── biz/
│       ├── handle/          # 【脚手架自动生成】HTTP接口实现层（核心业务入口）
│       ├── model/           # 【脚手架自动生成】请求/响应结构体、参数模型
│       ├── router/          # 【脚手架自动生成】路由注册代码
│       ├── mw/              # 手动维护：Hertz 中间件（鉴权、跨域、日志、限流等）
│       ├── logic/           # 手动维护：通用可复用业务逻辑层（与RPC逻辑层规范同源）
│       ├── convert/         # 手动维护：模型转换层（DB ↔ HTTP模型）
│       └── dal/             # 手动维护：数据访问层（与RPC DAL完全一致）
│           ├── gormL/
│           │   └── {db_name}/
│           │       ├── where/   # 链式查询构造器
│           │       └── *.go     # 表模型、Repo仓储
│           └── redis/
│               └── {db_name}/   # Redis数据层
├── render/                  # 页面渲染根目录（强制两层结构：业务模块/子页面）
│   ├── {biz_module}/        # 第一层：业务模块目录（如 shop、user）
│   │   ├── list.html        # 第二层：模块子页面（列表、新增、编辑、详情）
│   │   ├── add.html
│   │   ├── edit.html
│   │   └── detail.html
├── go.mod
├── go.sum
```

### 3\.1 核心目录权责说明

#### 3\.1\.1 自动生成目录（只读、禁止手动乱改）

**handle/、model/、router/**
1\. 由 Hertz 脚手架 `hz update` 命令根据 proto 文件全自动生成；
2\. **router/model 禁止手动修改**，迭代更新必须改 proto 重新生成；
3\. handle 层生成空方法模板，业务逻辑在模板内实现，不改动方法签名。

#### 3\.1\.2 idl 目录

1\. 存放所有 `.proto` 协议定义，统一管理 HTTP 接口、参数、结构体；
2\. 所有接口定义源头唯一，禁止代码内私自定义参数结构体。

#### 3\.1\.3 conf 目录

1\. 仅允许存放 **\.yaml 配置文件**；
2\. 禁止存放任何 Go 源码文件，严格做到配置与代码隔离。

#### 3\.1\.4 render 目录（Web 专属）

1\. 专门存放 HTML 页面模板、静态渲染资源，**目录强制固定两层业务层级**，禁止扁平化、禁止多层嵌套；
2\. 层级规范：第一层为**业务模块名**，第二层为**该模块下的子页面**（列表页、新增页、编辑页、详情页等）；
3\. 标准结构示例：

```Plain
render/
├── shop/                # 第一层：业务模块（店铺业务）
│   ├── list.html        # 第二层：子页面（列表页）
│   ├── add.html         # 第二层：子页面（新增页）
│   ├── edit.html        # 第二层：子页面（编辑页）
│   └── detail.html      # 第二层：子页面（详情页）
├── user/                # 第一层：业务模块（用户业务）
│   ├── list.html
│   ├── edit.html
│   └── add.html
```

4\. 统一使用 `github.com/osteele/liquid` 模板引擎渲染；
5\. 严格禁止：根目录直接存放 html、单业务多层嵌套、目录层级混乱；
6\. render 目录只存放静态模板与渲染资源，不包含任何业务 Go 代码。
1\. 手动维护，存放 Hertz 全局 / 路由中间件；
2\. 包含：跨域、鉴权、日志、限流、请求拦截、参数预处理等通用能力。

---

## 四、分层职责严格约束（Hertz Web 专属）

### 4\.1 分层核心职责

#### 4\.1\.1 biz/handle（HTTP 接口实现层）

1\. Hertz 接口唯一入口，接收前端 HTTP 请求；
2\. 负责参数校验、请求预处理、响应封装、异常捕获、日志打印；
3\. **允许编写当前接口专属轻量业务逻辑、缓存读写、参数组装**；
4\. 可复用、重型、通用业务逻辑禁止写在 handle，必须下沉至 logic；
5\. 禁止直接手写 SQL 操作数据库，必须通过 DAL/Repo 层操作。

#### 4\.1\.2 biz/logic（通用可复用业务层）

1\. 与 Kitex RPC 项目 logic 层规范完全同源；
2\. 存放**跨接口通用、可复用、重型业务逻辑、公共业务编排**；
3\. 统一收拢公共数据处理、复杂业务流程、聚合 DAL/Redis 能力；
4\. 无状态可全局复用，**禁止反向依赖 handle 上层**；
5\. 所有多接口复用逻辑必须下沉至此层，禁止在 handle 重复编写。

#### 4\.1\.3 biz/convert（模型转换层）

1\. 仅负责：**DB 模型 ↔ Hertz HTTP 模型** 双向字段映射、默认值补齐；
2\. 禁止编写任何业务判断、流程编排、数据处理逻辑；
3\. 职责绝对单一，只做模型转换。

#### 4\.1\.4 biz/dal（数据访问层）

1\. 与 RPC 项目 DAL 规范完全一致，多库隔离、结构统一；
2\. gormL：表模型、Repo 仓储、CURD 能力；
3\. where：链式查询构造器，类型安全、杜绝裸 SQL；
4\. redis：缓存读写逻辑；
5\. Repo 禁止全局单例，每次使用必须新建实例。

### 4\.2 分层依赖白名单（单向依赖、禁止反向）

\- biz/handle → logic、convert、dal、redis 缓存
\- biz/logic → convert、dal、redis 缓存
\- biz/convert → 无依赖
\- biz/dal → 无依赖

---

## 五、DAL 层强制规范（与 RPC 完全统一）

### 5\.1 Where 查询构造器规范

1\. 基类 `gormL.BaseWhere` 统一管理查询条件、SQL 构建；
2\. 所有查询必须使用链式 Where 构造，**禁止手写裸 SQL**；
3\. 文件名与表名对应小写，结构体为「大驼峰表名 \+ ListWhere」；
4\. **强制唯一入口**：禁止手动 new Where 结构体，必须通过对应 Repo 的 `GetWhere()` 方法获取实例。
**5\.1\.5 Where 构造器标准代码示例**

```Plain
package where
import "github.com/flyerxp/lib/v2/middleware/gormL"
// DemoListWhere 示例表查询构造器
type DemoListWhere struct {
	*gormL.BaseWhere
}
// GetDemoListWhere 统一创建Where实例
func GetDemoListWhere() *DemoListWhere {
	return &DemoListWhere{BaseWhere: &gormL.BaseWhere{}}
}
// TitleLike 标题模糊查询条件
func (w *DemoListWhere) TitleLike(title string) *DemoListWhere {
	w.Where("title LIKE ?", title+"%")
	return w
}
// StatusEq 状态精准匹配
func (w *DemoListWhere) StatusEq(status int) *DemoListWhere {
	w.Where("status = ?", status)
	return w
}
```

### 5\.2 Repo 仓储规范

1\. 每个 Repo 必须实现`GetWhere()` 方法，作为 Where 对象唯一获取入口；
2\. Repo 无状态、禁止全局单例，必须通过构造方法新建实例；
3\. 所有 DB 操作必须绑定 ctx 上下文；
4\. 严格遵循统一分页规范：limit\+1、自动判断 HasMore、自动裁剪列表；
5\. 接收外部事务`tx`参数时，必须对事务实例显式调用`Model`方法绑定当前 Repo 对应的表模型，禁止直接使用未绑定表模型的原始事务实例执行 DB 操作。

**5\.2\.1 Repo 仓储 \+ GORM 分页标准代码示例**

```Plain

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

**5\.2\.2 更新方法选型规范**

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

### 5\.3 全局分页规范

1\. 统一使用 `widget.Page` 分页结构体；
2\. limit\+1 查询方式精准判断下一页；
3\. 通过 `DoPage()` 统一裁剪数据、赋值分页状态。

### 5\.4 数据库初始化规范

`biz/dal/gormL/init.go` 为全局唯一数据库入口，多库统一管理、全局单次初始化，禁止分散连接。
**5\.4\.1 数据库初始化标准代码示例**

```Plain
package gormL
import (
	"context"
	"github.com/flyerxp/lib/v2/logger"
	"github.com/flyerxp/lib/v2/middleware/gormL"
	"gorm.io/gorm"
)
// DBClient 多库客户端聚合
type DBClient struct {
	Shop *gorm.DB
}
var dbClient *DBClient
// Init 全局数据库初始化
func Init(ctx context.Context) error {
	if dbClient != nil {
		return nil
	}
	shopDB, err := gormL.GetEngine(ctx, "shop")
	if err != nil {
		return err
	}
	dbClient = &DBClient{
		Shop: shopDB,
	}
	return nil
}
// GetShopDB 获取店铺库DB实例，自动绑定上下文+兜底初始化
func GetShopDB(ctx context.Context) *gorm.DB {
	if dbClient == nil {
		initCtx := logger.GetContext(context.Background(), "gormInit")
		if err := Init(initCtx); err != nil {
			panic(err)
		}
	}
	return dbClient.Shop.WithContext(ctx)
}
```

---

## 六、Hertz 脚手架生成规范

### 6\.1 生成原则

1\. 所有 HTTP 路由、模型、接口模板代码**必须通过 hz 脚手架生成**；
2\. 禁止手写路由注册、请求结构体、接口方法签名；
3\. 迭代更新接口必须先更新 idl 下 proto 文件，再执行生成命令；
4\. 固定使用 protobuf 模式，统一协议体系。

### 6\.2 生成目录约束

1\. router、model：纯自动生成，禁止手动修改，会被覆盖；
2\. handle：生成方法空模板，仅内部实现业务逻辑，不改动方法签名；
3\. 新增 / 修改接口必须以 proto 定义为唯一标准。

---

## 七、Web 页面渲染规范

1\. 所有前端 HTML 模板、渲染资源统一收敛在 `render/` 目录；
2\. 固定使用`github.com/osteele/liquid` 模板引擎做页面渲染；
3\. render 目录只存放静态模板与渲染资源，不包含任何业务 Go 代码；
4\. 页面数据组装、参数渲染逻辑统一放在 handle/logic 层，模板只负责展示。

---

## 八、全局红线禁止行为（MUST NOT）

1\. 禁止 logic 层反向依赖 handle 上层模块，禁止跨层、循环依赖；
2\. 禁止 convert 层编写任何业务判断、数据处理、流程逻辑；
3\. 禁止 handle 层编写可复用重型通用业务逻辑（必须下沉 logic）；
4\. 禁止手动 new 实例化 Where 结构体，必须通过 Repo\.GetWhere 获取；
5\. 禁止手写裸 SQL，所有查询必须使用 where 链式构造器；
6\. 禁止 Repo 全局单例常驻，必须每次新建实例；
7\. 禁止 conf 目录存放任何 Go 源码文件，只允许 yaml 配置；
8\. 禁止手动修改 router、model 自动生成代码；
9\. 禁止硬编码业务状态、时间、常量、ID 等数据；
10\. 禁止跨层、反向非法依赖，严格遵守单向分层白名单；
11\. 禁止 render 目录层级混乱，必须严格遵循「一级业务模块、二级子页面」两层结构；
12\. 禁止事务链式操作（Updates/Where等）不绑定表模型，严禁混用事务规范写法。

> （注：文档部分内容可能由 AI 生成）
> 
> 

> （注：文档部分内容可能由 AI 生成）
