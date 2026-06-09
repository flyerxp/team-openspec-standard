# OpenSpec 异步任务项目全局规范

## 优先级：全局最高，所有异步任务项目强制生效

一句话总结：project.md = 代码怎么写、工程怎么建

本规范适用于基于 Golang 构建的异步任务服务体系，包含定时任务（crontab）、消息消费（consumer）两类无 RPC 后台服务，统一工程结构、编码、数据访问、分页、分层依赖等核心规范，作为团队异步类项目统一工程模板。

## 一、项目核心信息

### 核心技术栈

Golang + GORM + Redis + Cobra + MQ + 定时任务（crontab）

### 规范基准

完全对齐 OpenSpec 全局 AI 代理规则、Kitex RPC DAL 全套规范



***

## 二、全局强制规约（MUST）

### 2.1 语言规范



1. 所有文档、注释、需求、设计方案、业务描述等必须为纯中文；

2. 禁止拼音、中英文混杂、无意义英文缩写；

3. 仅保留技术关键字英文：ctx、db、gorm、redis、cobra、req、resp；

4. 所有业务语义、变量释义、注释说明必须使用中文。

### 2.2 编码通用规范



1. 结构体统一聚合初始化：`var x T` / `x = T{}`；

2. 变量优先使用 `:=` 自动推导，禁止冗余类型声明；

3. 禁止变量遮蔽，内外作用域禁止重名变量；

4. 无限循环统一 `for{}`，禁止 `for 1`；

5. JSON、含双引号字符串必须使用反引号原生字符串；

6. 入口文件、启动文件禁止全局业务状态变量，状态全部结构体私有化；

7. Lint 可忽略告警行尾统一添加 `// NOLINT`；

8. 所有注释统一使用 `//` 单行注释；

9. Cobra 子命令统一规范：参数注册、默认值兜底、参数格式校验、标准错误输出；

10. service 目录与 cmd 子命令强一一对应，禁止目录混乱、代码散落。



***

## 三、全局固定目录结构（强制）

### 3.1 核心业务层 biz（对齐 Kitex）



```
biz/

├── logic            // 通用可复用业务逻辑层

├── convert          // 模型转换层

└── dal/

&#x20;   ├── gormL/

&#x20;   │   └── {db\_name}/

&#x20;   │       └── where/  // 链式查询构造器

&#x20;   └── redis/

&#x20;       └── {db\_name}/  // Redis数据层
```

### 3.2 定时任务 crontab



```
crontab/

├── cmd        // 命令行入口层

├── service    // 定时业务实现层

├── conf       // 定时专属配置

└── shell      // 调度脚本
```

### 3.3 消息消费 consumer



```
consumer/

├── cmd        // 命令行入口层

├── service    // 消费业务实现层

├── conf       // 消费专属配置

└── shell      // 启停执行脚本
```



***

## 四、分层职责严格约束



| 分层 / 模块                   | 核心职责                                                              | 禁止行为                              |
| ------------------------- | ----------------------------------------------------------------- | --------------------------------- |
| cmd（crontab/consumer）     | 命令注册、参数解析 / 校验、启动调度、资源释放                                          | 编写任何业务逻辑、数据查询、数据处理                |
| service（crontab/consumer） | 承载具体定时 / 消费业务编排与数据处理；统一处理上下文、日志、异常、分页批量逻辑；可调用 Logic、DAL、Convert   | 无（仅遵循依赖白名单）                       |
| Logic                     | 存放全局复用业务逻辑、通用工具、缓存逻辑；无状态、支持全局单例                                   | 依赖上层 service/cmd/convert；写 DAL 数据 |
| DAL                       | 按数据库分目录管理（多库隔离）；gormL 存放模型 / Repo/CURD、where 存放查询构造器、redis 存放缓存读写 | Repo 全局单例；手写裸 SQL；跨库非法调用          |
| Convert                   | 仅做 DB 模型 ↔ 业务模型双向转换                                               | 业务判断、流程编排、数据处理逻辑                  |



***

## 五、分层依赖白名单（强制）



* crontab/cmd → crontab/service

* consumer/cmd → consumer/service

* crontab/service → Logic、DAL、Convert

* consumer/service → Logic、DAL、Convert

* Logic → DAL（只读）、Convert



***

## 六、DAL 层完整强制规范

### 6.1 查询构造器规范



* 基类 BaseWhere 统一存储条件、Build 构建 SQL；

* 子类仅实现单行条件方法，职责单一，内部方法私有；

* 全程链式调用、类型安全、禁止裸 SQL；

#### 6.1.1 Where 命名规则



1. 文件名：`表名.go`，与 gormL 表模型文件一一对应；

2. 结构体名：`大驼峰表名 + ListWhere`。

#### 6.1.2 标准模板



```
package where

import "github.com/flyerxp/lib/v2/middleware/gormL"

type DemoListWhere struct {

&#x20;       \*gormL.BaseWhere

}

func (w \*DemoListWhere) TitleLike(title string) \*DemoListWhere {

&#x20;       w.Where("title LIKE ?", title+"%")

&#x20;       return w

}
```

### 6.2 Repo 强制规范



1. 所有 Repo 必须实现 `GetWhere()` 方法，返回对应 Where 实例；

2. Repo 禁止全局单例，必须提供构造方法每次新建实例。

#### 6.2.1 标准 Repo + 分页代码示例

```

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
}
// ModuleInfo 模块示例表模型（适配Save事务示例）
type ModuleInfo struct {
	Id     int    `gorm:"column:id;primaryKey;autoIncrement"`
	Path   string `gorm:"column:path;size:500"`
	RootId int    `gorm:"column:root_id"`
}
// TableName 指定表名
func (DemoInfo) TableName() string {
	return "demo_info"
}
// TableName 指定表名
func (ModuleInfo) TableName() string {
	return "module_info"
}
// DemoRepo 数据仓储层
type DemoRepo struct{}
// ModuleRepo 模块数据仓储层
type ModuleRepo struct{}

// NewDemoRepo 构造方法：每次新建实例，禁止全局单例
func NewDemoRepo() *DemoRepo {
	return &DemoRepo{}
}
// NewModuleRepo 构造方法：每次新建实例，禁止全局单例
func NewModuleRepo() *ModuleRepo {
	return &ModuleRepo{}
}

// GetWhere 强制实现：Where对象唯一获取入口
func (r *DemoRepo) GetWhere() *where.DemoListWhere {
	return &where.DemoListWhere{BaseWhere: &gormLib.BaseWhere{}} 
}

// GetGormModel 获取绑定上下文的DB实例
func (r *DemoRepo) GetGormModel(ctx context.Context) *gorm.DB {
	return gormL.GetDB(ctx).Model(&DemoInfo{})
}

// GetGormModel 获取绑定上下文的DB实例
func (r *ModuleRepo) GetGormModel(ctx context.Context) *gorm.DB {
	return gormL.GetDB(ctx).Model(&ModuleInfo{})
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
	db := r.GetGormModel(ctx)
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

// UpdatePathById 事务/非事务通用更新方法（链式Updates操作规范示例）
// 适用场景：Updates、Where、Select等链式更新操作
func (r *DemoRepo) UpdatePathById(ctx context.Context, id int, path string, rootId int, tx *gorm.DB) error {
	db := tx
	if db == nil {
		// 无事务：使用预设表模型的标准DB实例
		db = r.GetGormModel(ctx)
	} else {
		// 【强制红线】事务链式操作必须手动绑定当前表模型
		db = db.Model(&DemoInfo{})
	}
	return db.Where("id = ?", id).Updates(map[string]interface{}{
		"path":    path,
		"root_id": rootId,
	}).Error
}

// UpdatePathById 事务/非事务通用保存方法（Save结构体保存规范示例）
// 适用场景：完整模型结构体保存、更新，依托结构体自动识别表模型
func (r *ModuleRepo) UpdatePathById(ctx context.Context, info *ModuleInfo, tx *gorm.DB) error {
	// 事务场景：Save依托结构体自动识别表模型，无需手动Model绑定
	if tx != nil {
		return tx.Save(info).Error
	}
	// 非事务场景：统一使用标准绑定模型DB实例
	return r.GetGormModel(ctx).Save(info).Error
}
```

#### 6.2.2 事务参数处理强制规范

当业务需要跨 Repo 执行事务操作，向 Repo 方法传递外部事务实例时，必须严格遵循以下规则：

1. **链式操作（Updates/Where 等）强制绑定表模型**：外部传入的事务`*gorm.DB`为通用实例，未绑定当前Repo表模型，执行 `Updates`、`Where`、`Select` 等链式操作时，**必须手动调用 Model\(\) 绑定当前业务表模型**，否则会出现表名识别错误、字段映射失效、软删除/钩子不生效等问题。

2. **结构体保存（Save）豁免绑定 Model**：使用 GORM `Save` 方法时，方法入参为完整模型结构体，GORM 可自动识别表模型，**无需手动调用 Model\(\) 绑定**，但必须严格做事务非空判断，统一事务、非事务分支代码逻辑，保持风格一致。

3. **事务分支必做判空兜底**：所有含事务参数的Repo方法，必须区分「有事务」「无事务」两个分支，无事务时统一使用 `GetGormModel` 获取已绑定表模型的标准DB实例，禁止直接裸用DB。

4. **统一事务行为、杜绝隐性BUG**：所有事务操作严格区分两类写法适配规则，不混用、不省略绑定逻辑，保证所有数据库操作行为统一，规避事务场景下的数据更新异常、SQL 执行错误、事务回滚失效等隐性问题。

### 6.3 统一分页规范



* 统一分页结构体，自动计算 `HasMore`；

* 统一 `limit+1` 查询方式，精准判断是否有下一页；

* 统一 `DoPage()` 自动裁剪列表数据；

* 所有批量查询、遍历查询必须遵守该规范。

### 6.4 gormL 数据库初始化规范

biz/dal/gormL/init.go 为项目全局数据库入口文件，负责多数据库实例初始化、单例管理、上下文绑定，全局唯一，禁止分散初始化。

#### 标准 init.go 示例



```
package gormL

import (

&#x20;       "context"

&#x20;       "github.com/flyerxp/lib/v2/logger"

&#x20;       "github.com/flyerxp/lib/v2/middleware/gormL"

&#x20;       "gorm.io/gorm"

)

type DBClient struct {

&#x20;       Shop   \*gorm.DB

&#x20;       Report \*gorm.DB

}

var dbClient \*DBClient

func Init(ctx context.Context) error {

&#x20;       if dbClient != nil {

&#x20;               return nil

&#x20;       }

&#x20;       shopDB, err := gormL.GetEngine(ctx, "readshop")

&#x20;       if err != nil {

&#x20;               return err

&#x20;       }

&#x20;       reportDB, err := gormL.GetEngine(ctx, "report")

&#x20;       if err != nil {

&#x20;               return err

&#x20;       }

&#x20;       dbClient = \&DBClient{Shop: shopDB, Report: reportDB}

&#x20;       return nil

}

func GetShopDB(ctx context.Context) \*gorm.DB {

&#x20;       if dbClient == nil {

&#x20;               initCtx := logger.GetContext(context.Background(), "gormInit")

&#x20;               if err := Init(initCtx); err != nil {

&#x20;                       panic(err)

&#x20;               }

&#x20;       }

&#x20;       return dbClient.Shop.WithContext(ctx)

}

func GetReportDB(ctx context.Context) \*gorm.DB {

&#x20;       if dbClient == nil {

&#x20;               initCtx := logger.GetContext(context.Background(), "gormInit")

&#x20;               if err := Init(initCtx); err != nil {

&#x20;                       panic(err)

&#x20;               }

&#x20;       }

&#x20;       return dbClient.Report.WithContext(ctx)

}
```



***

## 七、项目红线禁止行为（MUST）



1. 禁止 Logic 依赖上层 crontab/consumer 模块；

2. 禁止 cmd 层编写任何业务逻辑、数据操作；

3. 禁止 Convert 编写流程编排、数据处理逻辑；

4. 禁止 DAL Repo 全局单例常驻，必须每次新建实例；

5. 禁止硬编码配置、时间、状态值、ID 等业务常量；

6. 禁止手写裸 SQL，全部使用 where 链式构造器；

7. 禁止跨层非法依赖、反向依赖。