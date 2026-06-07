# AGENTS\.md

# OpenSpec 全局 AI 代理规则（定时/消费专属项目 · 对齐Kitex标准）

## 优先级：全局最高，所有项目强制生效

## 1\. 语言强制规则（MUST）

1\. 所有文档、注释、需求、设计、方案必须纯中文

2\. 禁止拼音、禁止中英文混杂、禁止无意义英文缩写

3\. 技术关键字保留英文：ctx、db、req、resp、gorm、redis、cobra

4\. 业务语义全部使用中文描述

## 2\. 框架固定架构（MUST）

技术栈：Golang \+ GORM \+ Redis \+ 定时任务 \+ MQ消息消费 \+ cobra命令行

分层固定：Logic / DAL / Convert

上层调度模块：crontab（定时任务）、consumer（消息消费）

### 分层职责严格约束

#### Logic（通用业务逻辑层）

\- 存放项目所有可复用业务逻辑、缓存逻辑、工具方法、业务编排

\- 承接crontab、consumer所有通用核心业务能力调用

\- 无状态、可全局单例

\- 禁止依赖任何上层调用模块（crontab/consumer）、Convert

\- 仅可只读读取 DAL 数据，可调用Redis缓存逻辑

#### DAL（数据访问层）

\- 完全对齐团队Kitex RPC项目规范，严格按数据库名称分目录管理

\- 每个库独立：gormL / gormL/where / redis

\- gormL：模型定义、Repo仓储、通用CURD能力

\- where：仿 Ent 链式查询构造器、基类继承、私有方法封装

\- Repo 无状态，每次 new 新对象，不使用单例

\- 支持被 logic、crontab/service、consumer/service 直接调用

#### Convert（模型转换层）

\- 只做 DB模型 ↔ 业务模型 双向转换

\- 禁止业务逻辑、禁止判断、禁止流程编排

#### crontab/consumer 上层模块职责

##### cmd 命令行入口层

\- 基于 spf13/cobra 实现命令行注册、参数接收、参数校验、命令分发

\- 仅做参数解析、命令绑定、入口调度、资源释放，**禁止编写任何业务逻辑**

\- 统一注册根命令、子命令、自定义参数、默认参数兜底

##### service 业务实现层

\- 对应cmd子命令，按命令维度拆分目录与业务文件，目录文件名一一对应

\- 实现具体定时任务、消费任务核心业务流程与数据处理

\- 可直接调用 biz/logic 通用能力、直接调用 biz/dal 数据层CURD

\- 统一处理上下文、日志、异常、告警、批量分页处理逻辑

##### conf 配置层

\- 存放当前模块专属配置文件，自动生成后人工微调参数

\- 禁止硬编码配置，所有动态参数、业务配置统一收敛至配置文件

##### shell 脚本层

\- 存放批量执行、定时调度、权限处理、服务启停脚本

\- 统一处理日志权限、执行用户、批量任务调度

## 3\. 目录结构强制固定（完全对齐Kitex标准）

biz/dal/gormL/\{db\_name\}

biz/dal/gormL/\{db\_name\}/where

biz/dal/redis/\{db\_name\}

**上层模块固定四层目录**

crontab/cmd

crontab/service

crontab/conf

crontab/shell

consumer/cmd

consumer/service

consumer/conf

consumer/shell

## 4\. Golang 编码强制规范（MUST）

1\. 结构体必须聚合初始化：var x T / x = T\{\}

2\. 变量赋值优先 := 自动推导，禁止冗余类型声明

3\. 禁止变量遮蔽，内部作用域禁止与外部重名

4\. 无限循环统一 for\{\}，禁止 for 1

5\. JSON/含双引号字符串，必须使用反引号原生字符串

6\. 所有cmd、main启动文件禁止定义全局业务状态变量，所有状态封装结构体私有变量

7\. Lint 可静态化警告，行尾加 // NOLINT

8\. 所有注释使用 //，禁止其他注释方式

9\. cobra子命令统一规范：参数注册、默认值兜底、参数格式校验、标准错误输出

10\. service目录与cmd子命令强一一对应，禁止目录混乱、代码散落

## 5\. 查询构造器固定规范（与Kitex项目完全统一）

\- 基类 BaseWhere 统一存储条件、Build 构建SQL能力

\- 子类只写单行 where 条件逻辑，职责单一

\- 内部方法私有，外部无法调用内部私有方法

\- 全程链式调用、类型安全、无裸 SQL

\- 完全对齐团队Kitex RPC项目查询构造器标准

### 5\.1 Where 文件命名规则

1\. where目录文件名：`表名.go`，和gormL目录下表定义文件同名一一对应；

2\. where内部结构体命名：`表名大驼峰+ListWhere`。

#### Where标准代码模板

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

### 5\.2 gormL Repo 强制规范

1\. 每个Repo必须自带GetWhere\(\)方法，统一获取对应Where实例；

2\. Repo禁止全局单例，通过构造方法新建实例，用完即释。

#### 标准Repo代码示例：

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

## 6\. 分页规范（全局统一，对齐Kitex标准）

\- 统一分页结构体、自动判断 HasMore 有无下一页

\- 统一 limit\+1 分页查询逻辑，精准判断分页边界

\- 统一列表裁剪逻辑，自动截断多余数据

\- 所有批量查询、列表查询必须遵循该分页规范

## 7\. 分层依赖白名单（强制）

crontab/cmd    → crontab/service

consumer/cmd   → consumer/service

crontab/service   → Logic、DAL、Convert

consumer/service  → Logic、DAL、Convert

Logic → DAL\(只读\)、Convert

## 8\. 禁止行为红线（MUST）

\- 禁止 Logic 依赖 crontab、consumer 上层模块

\- 禁止 cmd 层编写任何业务逻辑、数据操作

\- 禁止 Convert 编写任何业务逻辑、判断逻辑

\- 禁止 DAL Repo 全局单例常驻，必须每次新建实例

\- 禁止业务代码硬编码配置、时间、ID、状态值

\- 禁止手写裸SQL，所有查询必须通过where链式构造器实现

\- 禁止跨层非法依赖、反向依赖

> （注：文档部分内容可能由 AI 生成）
