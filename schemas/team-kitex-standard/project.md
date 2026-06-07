# project\.md（异步任务项目工程规范）

## 一、项目简介

本项目为 Golang 异步任务服务体系，包含 **定时任务（crontab）**、**消息消费（consumer）** 两类无RPC后台服务，基于 Cobra 命令行驱动，完全对齐团队 Kitex 微服务 DAL 层标准，统一工程结构、编码、查询构造器、分页、分层依赖规范，为团队异步类项目统一工程模板。

**技术栈**：Golang \+ GORM \+ Redis \+ Cobra \+ MQ

**规范基准**：完全对齐 OpenSpec 全局 AI 代理规则、Kitex RPC DAL 全套规范

## 二、全局强制规约（MUST）

### 2\.1 语言规范

1\. 所有注释、文档、业务描述、设计方案 **纯中文**。

2\. 禁止拼音、中英文混杂、无意义英文缩写。

3\. 仅保留技术关键字英文：ctx、db、gorm、redis、cobra。

4\. 所有业务语义、变量释义、注释说明必须中文。

### 2\.2 编码通用规范

1\. 结构体统一聚合初始化：`var x T` / `x = T{}`。

2\. 变量优先使用 `:=` 自动推导，禁止冗余类型声明。

3\. 禁止变量遮蔽，内外作用域禁止重名变量。

4\. 无限循环统一 `for{}`，禁止 `for 1`。

5\. JSON、含双引号字符串必须使用反引号原生字符串。

6\. 入口文件、启动文件禁止全局业务状态变量，状态全部结构体私有化。

7\. Lint 可忽略告警行尾统一添加 `// NOLINT`。

8\. 所有注释统一使用 `//` 单行注释。

## 三、全局固定目录结构（强制）

### 3\.1 核心 biz 目录（对齐Kitex）

biz/logic          通用可复用业务逻辑层

biz/convert        模型转换层（仅做映射）

biz/dal/gormL/\{db\_name\}    数据库模型 \& Repo层

biz/dal/gormL/\{db\_name\}/where  链式查询构造器

biz/dal/redis/\{db\_name\}    Redis数据层

### 3\.2 定时任务 crontab 四层目录

crontab/cmd     命令行入口层

crontab/service 定时业务实现层

crontab/conf    定时专属配置

crontab/shell   调度脚本

### 3\.3 消息消费 consumer 四层目录

consumer/cmd    命令行入口层

consumer/service 消费业务实现层

consumer/conf   消费专属配置

consumer/shell  启停执行脚本

## 四、分层职责严格约束

### 4\.1 cmd 层

\- 仅负责：命令注册、参数解析、参数校验、启动调度、资源释放

\- 禁止：任何业务逻辑、数据查询、数据处理

### 4\.2 service 层（crontab/consumer）

\- 承载具体定时、消费业务编排与数据处理

\- 可直接调用 Logic、DAL、Convert

\- 统一处理上下文、日志、异常、分页批量逻辑

### 4\.3 Logic 层

\- 存放全局复用业务逻辑、通用工具、缓存逻辑

\- 无状态、支持全局单例

\- 禁止依赖上层 service、cmd、convert

\- 仅可只读读取 DAL 数据

### 4\.4 DAL 层

\- 严格按数据库分目录管理，多库完全隔离

\- gormL 存放模型、Repo、CURD逻辑

\- where 存放链式查询构造器

\- redis 存放缓存读写逻辑

\- Repo 无状态，每次使用必须 new 新实例，禁止单例

### 4\.5 Convert 层

\- 只做：DB模型 ↔ 业务模型 双向转换

\- 禁止：业务判断、流程编排、数据处理逻辑

## 五、DAL 层完整强制规范

### 5\.1 查询构造器规范

\- 基类 BaseWhere 统一存储条件、Build 构建SQL

\- 子类仅实现单行条件方法，职责单一

\- 内部方法私有，外部不可调用

\- 全程链式调用、类型安全、禁止裸SQL

\- 与 Kitex RPC 项目完全一致

#### 5\.1\.1 Where 命名规则

1\. 文件名：`表名.go`，与 gormL 表模型文件一一对应

2\. 结构体名：`大驼峰表名 + ListWhere`

#### 5\.1\.2 Where 标准模板

```go
package where

import "github.com/flyerxp/lib/v2/middleware/gormL"

type NewsListWhere struct {
	*gormL.BaseWhere
}

func (w *NewsListWhere) TitleLike(title string) *NewsListWhere {
	w.Where("title LIKE ?", title+"%")
	return w
}

```

### 5\.2 Repo 强制规范

1\. 所有 Repo 必须实现 `GetWhere()` 方法，返回对应Where实例

2\. Repo 禁止全局单例，必须提供构造方法每次新建实例

#### 5\.2\.1 标准 Repo \+ 分页代码示例

```go
package ch123

import (
	"context"
	"github.com/flyerxp/lib/v2/middleware/gormL"
	"github.com/flyerxp/content.news.rpc/v2/biz/dal/gormL/ch123/where"
	"github.com/flyerxp/globalStruct/widget"
	"gorm.io/gorm"
)

// NewsInfo 新闻资讯 GORM 模型
type NewsInfo struct {
	Id           int       `gorm:"column:id;primaryKey;autoIncrement" json:"id,omitempty"`
	Title        string    `gorm:"column:title;size:255" json:"id,omitempty"`
	CategoryId   int       `gorm:"column:category_id" json:"category_id,omitempty"`
}
func (NewsInfo) TableName() string {return "news_info"}

// NewsRepo 数据仓储
type NewsRepo struct{}
// GetNewNewsRepo 新建实例
func GetNewNewsRepo() *NewsRepo {return &NewsRepo{}}
// GetWhere 强制方法，获取对应Where对象
func (n *NewsRepo) GetWhere() *where.NewsListWhere {
	return &where.NewsListWhere{BaseWhere: &gormLib.BaseWhere{}}
}
func (r *NewsRepo) getDB(ctx context.Context) *gorm.DB {
	return gormL.GetDB(ctx).Model(&NewsInfo{})
}
// NewsInfoListColsPage 分页结构体
type NewsInfoListColsPage struct {
	List []NewsListCols
	Page widget.Page
}
func (n *NewsInfoListColsPage) DoPage() *NewsInfoListColsPage {
	n.Page.HasMore = len(n.List) > n.Page.Size
	if n.Page.HasMore {
		n.List = n.List[:n.Page.Size]
	}
	return n
}

// GetList 分页查询
func (r *NewsRepo) GetList(ctx context.Context, w *where.NewsListWhere, sort string, page int, limit int) (*NewsInfoListColsPage, error) {
	var list []NewsListCols
	pageObj := NewsInfoListColsPage{List: list, Page: widget.Page{Size: limit,Page: page}}
	db := r.GetDB(ctx).Select("id","title","description","create_time","update_time","subtitle","time_line","img","category_id","status")
	if w != nil {
		db = w.Build(db)
	}
	switch sort {
	case "web":
		db = db.Order("is_top desc, sort_id desc, update_time desc")
	case "time_line":
		db = db.Order("is_top desc, sort_id desc, time_line desc")
	default:
		db = db.Order("id desc")
	}
	offset := (page - 1) * limit
	db = db.Offset(offset).Limit(limit + 1)
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	pageObj.List = list
	return pageObj.DoPage(), nil
}

```

### 5\.3 统一分页规范

\- 统一分页结构体，自动计算 `HasMore`

\- 统一 `limit+1` 查询方式，精准判断是否有下一页

\- 统一 `DoPage()` 自动裁剪列表数据

\- 所有批量查询、遍历查询必须遵守该规范

### 5\.3 统一分页规范

\- 统一分页结构体，自动计算 `HasMore`

\- 统一 `limit+1` 查询方式，精准判断是否有下一页

\- 统一 `DoPage()`自动裁剪列表数据

\- 所有批量查询、遍历查询必须遵守该规范

### 5\.4 gormL 数据库初始化规范（init\.go）

biz/dal/gormL/init\.go 为项目全局数据库入口文件，负责多数据库实例初始化、单例管理、上下文绑定，是所有DAL数据查询的统一入口，全局唯一，禁止分散初始化数据库连接。

**核心约束**：

1\. 支持多数据库并行初始化，统一聚合到 DBClient 结构体管理

2\. 全局单次初始化，避免重复创建数据库连接

3\. 所有DB操作必须绑定上下文 ctx，禁止裸DB实例操作

4\. 内部兜底自动初始化，无需业务层手动调用 Init 方法

5\. 按业务库拆分独立获取方法，职责隔离、清晰可控

**标准 init\.go 完整示例**：

```go
package gormL

import (
	"context"
	"github.com/flyerxp/lib/v2/logger"
	"github.com/flyerxp/lib/v2/middleware/gormL"
	"gorm.io/gorm"
)

// DBClient 定义GORM数据库连接聚合实例
type DBClient struct {
	Shop   *gorm.DB
	Report *gorm.DB
}

var dbClient *DBClient

// Init 初始化GORM多数据库连接（全局仅执行一次）
func Init(ctx context.Context) error {
	if dbClient != nil {
		return nil
	}
	var shop *gorm.DB
	var report *gorm.DB
	var err error
	if shop, err = gormL.GetEngine(ctx, "readshop"); err != nil {
		return err
	}

	report, err = gormL.GetEngine(ctx, "report")
	if err != nil {
		return err
	}
	dbClient = &DBClient{Shop: shop, Report: report}
	return nil
}

// GetShopDB 获取已注入context的Shop库DB实例
func GetShopDB(ctx context.Context) *gorm.DB {
	if dbClient == nil {
		initCtx := logger.GetContext(context.Background(), "gormInit")
		if err := Init(initCtx); err != nil {
			panic(err)
		}
	}
	return dbClient.Shop.WithContext(ctx)
}

// GetReportDB 获取已注入context的Report库DB实例
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

## 六、分层依赖白名单（强制）

crontab/cmd    → crontab/service

consumer/cmd   → consumer/service

crontab/service   → Logic、DAL、Convert

consumer/service  → Logic、DAL、Convert

Logic → DAL\(只读\)、Convert

## 七、项目红线禁止行为

\- 禁止 Logic 依赖上层 crontab/consumer 模块

\- 禁止 cmd 层写任何业务、数据逻辑

\- 禁止 Convert 写判断、转换以外的逻辑

\- 禁止 DAL Repo 全局单例常驻

\- 禁止硬编码配置、时间、状态值

\- 禁止手写裸SQL，全部使用where链式构造器

\- 禁止跨层非法依赖、反向依赖

> （注：文档部分内容可能由 AI 生成）
