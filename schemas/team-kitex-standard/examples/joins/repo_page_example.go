// 文件路径：biz/dal/gormL/shop/joins/demo_info.go
package joins

import (
	"context"
	rpcExamples "rpcExamples/examples"
	"time"

	"github.com/flyerxp/globalStruct/widget"
	"gorm.io/gorm"

	"rpcExamples/examples/joins/where"
)

// DemoInfoJoinRow 连表查询结果承接DTO
// 严格合规：无TableName、无gorm.Model、无gorm关联标签、仅column映射
type DemoInfoJoinRow struct {
	Id          int       `gorm:"column:id"`
	Title       string    `gorm:"column:title"`
	Status      int       `gorm:"column:status"`
	Path        string    `gorm:"column:path"`
	RootId      int       `gorm:"column:root_id"`
	PriceId     int       `gorm:"column:price_id"`
	StrategyId  int       `gorm:"column:strategy_id"`
	PriceAmount float64   `gorm:"column:price_amount"`
	CreateTime  time.Time `gorm:"column:create_time"`
}

// DemoInfoSimpleJoinRow 精简JSON列表承接DTO
type DemoInfoSimpleJoinRow struct {
	Id          int     `gorm:"column:id"`
	Title       string  `gorm:"column:title"`
	PriceAmount float64 `gorm:"column:price_amount"`
}

// DemoInfoJoinPage 统一分页返回结构体
type DemoInfoJoinPage struct {
	List []DemoInfoJoinRow
	Page widget.Page
}

// DoPage 分页裁剪、自动计算是否有下一页
func (p *DemoInfoJoinPage) DoPage() *DemoInfoJoinPage {
	p.Page.HasMore = len(p.List) > p.Page.Size
	if p.Page.HasMore {
		p.List = p.List[:p.Page.Size]
	}
	return p
}

// DemoInfoJoinRepo 连表仓储（demo_info + demo_info_price INNER/LEFT JOIN）
type DemoInfoJoinRepo struct{}

// NewDemoInfoJoinRepo 构造方法：每次新建实例，禁止全局单例
func NewDemoInfoJoinRepo() *DemoInfoJoinRepo {
	return &DemoInfoJoinRepo{}
}

// GetWhere 全局统一命名获取Where构造器，对齐user模板
func (r *DemoInfoJoinRepo) GetWhere() *where.DemoJoinWhere {
	return where.NewDemoJoinWhere()
}

// getGormModel 抽离DB初始化：直接Table绑定原生表名，内部内置JOIN语句
// 不再依赖单表repo，完全对齐user模板范式
func (r *DemoInfoJoinRepo) getGormModel(ctx context.Context, isInner bool) *gorm.DB {
	db := rpcExamples.GetShopDB(ctx).Table("demo_info")
	// 动态切换内外连接，和user模板保持一致扩展逻辑
	if isInner {
		db = db.Joins("INNER JOIN `demo_info_price` ON `demo_info_price`.`strategy_id` = `demo_info`.`root_id`")
	} else {
		db = db.Joins("LEFT JOIN `demo_info_price` ON `demo_info_price`.`strategy_id` = `demo_info`.`root_id`")
	}
	return db
}

// ListPage 分页查询（对应业务列表分页）
func (r *DemoInfoJoinRepo) ListPage(ctx context.Context, w *where.DemoJoinWhere, page, limit int, isInner bool) (*DemoInfoJoinPage, error) {
	var list []DemoInfoJoinRow
	pageObj := DemoInfoJoinPage{
		List: list,
		Page: widget.Page{Page: page, Size: limit},
	}

	db := r.getGormModel(ctx, isInner)
	// 显式指定所有查询字段，严格禁止SELECT *，字段全部携带原生表名
	db = db.Select([]string{
		"`demo_info`.`id`",
		"`demo_info`.`title`",
		"`demo_info`.`status`",
		"`demo_info`.`path`",
		"`demo_info`.`root_id`",
		"`demo_info_price`.`id` AS price_id",
		"`demo_info_price`.`strategy_id`",
		"`demo_info_price`.`price_amount`",
		"`demo_info`.`create_time`",
	})
	// 统一硬删除过滤，所有业务表默认逻辑
	db = db.Where("`demo_info`.`is_delete` = 0")
	// 拼接连表查询条件
	if w != nil {
		db = w.Build(db)
	}
	// 排序必须携带原生表名，杜绝字段歧义
	db = db.Order("`demo_info`.`id` DESC")
	// limit+1分页判断更多
	offset := (page - 1) * limit
	db = db.Offset(offset).Limit(limit + 1)

	if err := db.Find(&list).Error; err != nil {
		return pageObj.DoPage(), err
	}
	pageObj.List = list
	return pageObj.DoPage(), nil
}

// ListSimple 精简列表查询（对应前端下拉/JSON简单列表）
func (r *DemoInfoJoinRepo) ListSimple(ctx context.Context, w *where.DemoJoinWhere, limit int, isInner bool) ([]DemoInfoSimpleJoinRow, error) {
	var list []DemoInfoSimpleJoinRow

	db := r.getGormModel(ctx, isInner)
	db = db.Select([]string{
		"`demo_info`.`id`",
		"`demo_info`.`title`",
		"`demo_info_price`.`price_amount`",
	})
	db = db.Where("`demo_info`.`is_delete` = 0")
	if w != nil {
		db = w.Build(db)
	}
	db = db.Order("`demo_info`.`id` DESC")
	db = db.Limit(limit + 1)

	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
