# OpenSpec 全局 AI 代理规则（团队 Hertz Web 服务专用）
## 优先级：全局最高，仅 Hertz 项目强制生效，不影响 Kitex RPC 项目
# 1. 语言强制规则（MUST）
1. 所有文档、注释、需求、设计、方案必须纯中文
2. 禁止拼音、禁止中英文混杂、禁止无意义英文缩写
3. 技术关键字保留英文：ctx、db、req、resp、gorm、hertz
4. 业务语义全部使用中文描述
5. 与 Kitex 项目语言规范完全统一，保证团队风格一致
# 2. 框架固定架构（MUST）
技术栈：Golang + Hertz Web + GORM + Redis
分层固定：Controller / Service / Logic / DAL / Convert
区别于Kitex RPC架构，严格遵循Web分层规范
## 分层职责严格约束
### Controller（接口控制器层）
- Web 请求唯一入口，对应路由接口
- 只做：参数校验绑定、请求透传、日志埋点、错误包装、响应返回
- 禁止：任何业务逻辑、数据查询、数据处理、事务操作
### Service（业务实现层）
- 一对一实现 Web 接口业务逻辑
- 负责业务流程编排、参数二次校验、数据组装、事务控制
- 允许直接操作 DB、Redis（通过 DAL Repo）
- 允许调用 Logic 通用方法、Convert 模型转换
### Logic（通用逻辑层）
- 存放多 Web 接口复用逻辑、缓存逻辑、通用工具、计算规则
- 无状态、可全局单例
- 禁止依赖 Service、Controller、Convert
- 仅可只读读取 DAL 数据
- 与 Kitex 项目 Logic 层规范完全一致
### DAL（数据访问层）
- 严格按数据库名称分目录：biz/dal/{db_name}/
- 每个库独立：gormL / gormL/where / redis
- gormL：模型、Repo、CURD
- where：仿 Ent 链式查询构造器、基类继承、私有方法封装
- Repo 无状态，每次 new 新对象，不使用单例
- 与 Kitex 项目 DAL 目录结构完全统一，保证数据层规范一致
### Convert（模型转换层）
- 只做 DB结构体 ↔ Web入参/出参结构体 转换
- 禁止业务逻辑、禁止判断、禁止流程编排
- 统一团队转换规范
# 3. Hertz 专属固定目录结构
biz/controller       # Web控制器
biz/dal/gormL/{db_name}
biz/dal/gormL/{db_name}/where
biz/dal/{db_name}/redis
router/              # 路由注册专属目录
# 4. Golang 编码强制规范（MUST，全团队统一）
1. 结构体必须聚合初始化：var x T / x = T{}
2. 变量赋值优先 := 自动推导，禁止冗余类型声明
3. 禁止变量遮蔽，内部作用域禁止与外部重名
4. 无限循环统一 for{}，禁止 for 1
5. JSON、含双引号字符串，必须使用反引号原生字符串
6. main.go 禁止定义全局状态变量，所有状态封装结构体私有变量
7. Lint 可静态化警告，行尾加 // NOLINT
8. 所有注释使用 //，禁止其他注释方式
# 5. 查询构造器固定规范
- 基类 BaseWhere 统一存储条件、Build 构建
- 子类只写一行 where 条件
- 外部无法调用内部私有方法
- 全程链式调用、类型安全、无裸 SQL
- 与 Kitex 项目完全统一

## 5.1 Where 文件命名规则
1. where目录文件名：`表名.go`，和gormL目录下表定义文件同名一一对应；
2. where内部结构体命名：`表名大驼峰+ListWhere`。

### Where标准代码模板
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
## 5.2 gormL Repo 强制规范
1. 每个Repo必须自带GetWhere()方法，统一获取对应Where实例；
2. Repo禁止全局单例，通过构造方法新建实例。

**标准Repo代码示例：**
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
	return &where.NewsListWhere{}
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
	db := gormL.GetDB(ctx).Model(&NewsInfo{}).Select("id","title","description","create_time","update_time","subtitle","time_line","img","category_id","status")
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
# 6. 分页规范
- 统一分页结构体、自动判断 HasMore
- 统一 limit+1 分页逻辑
- 统一列表裁剪逻辑

# 7. 禁止行为红线
- 禁止 Logic 依赖上层业务
- 禁止 Controller 写业务逻辑
- 禁止 Convert 写逻辑
- 禁止 DAL 单例常驻
- 禁止硬编码配置

# 8. 双框架隔离规则
- 本规则仅生效 Kitex RPC 项目
- 不干预、不污染 Hertz Web 项目架构
- 基础编码、语言规范双向统一，架构分层完全隔离