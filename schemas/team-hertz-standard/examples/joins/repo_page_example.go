// 文件路径：biz/dal/gormL/shop/joins/demo_info.go
package joins

import (
	"context"
	"github.com/flyerxp/globalStruct/widget"
	"gorm.io/gorm"
	webExamples "webExamples/examples"
	"webExamples/examples/joins/where"
)

// DemoInfoJoinRow 连表查询外层承接DTO 仅做结果映射，无任何关联标签
type DemoInfoJoinRow struct {
	Id     int    `gorm:"column:id"`
	Title  string `gorm:"column:title"`
	Status int    `gorm:"column:status"`
	Path   string `gorm:"column:path"`
	RootId int    `gorm:"column:root_id"`
	// 关联精简DTO，承接price_id别名字段
	//Price DemoInfoPriceSimple
	// 关联精简DTO，承接price_id别名字段,一对多
	Price []DemoInfoPriceSimple
}

// DemoInfoPriceSimple 关联表结果承接DTO 【核心注意】不识别表名、无TableName
type DemoInfoPriceSimple struct {
	PriceId     int     `gorm:"column:price_id"`
	StrategyId  int     `gorm:"column:strategy_id"`
	PriceAmount float64 `gorm:"column:price_amount"`
}

// DemoJoinPage 统一分页返回结构体
type DemoJoinPage struct {
	List []DemoInfoJoinRow
	Page widget.Page
}

// DoPage 分页溢出裁剪、是否还有更多
func (p *DemoJoinPage) DoPage() *DemoJoinPage {
	p.Page.HasMore = len(p.List) > p.Page.Size
	if p.Page.HasMore {
		p.List = p.List[:p.Page.Size]
	}

	return p
}

// DemoInfoJoinRepo 连表仓储
type DemoInfoJoinRepo struct{}

func NewDemoInfoJoinRepo() *DemoInfoJoinRepo {
	return &DemoInfoJoinRepo{}
}

// GetWhere 全局统一命名，和单表repo方法对齐，废弃原有GetJoinWhere
func (r *DemoInfoJoinRepo) GetWhere() *where.DemoJoinWhere {
	return where.NewDemoJoinWhere()
}
func (r *DemoInfoJoinRepo) getGormModel(ctx context.Context) *gorm.DB {

	// 无事务：使用默认连接
	return webExamples.GetShopDB(ctx).Model(&webExamples.DemoInfo{})
}

// ListJoinPage 【完整无删减】Preload连表分页查询
func (r *DemoInfoJoinRepo) ListJoinPage(ctx context.Context, w *where.DemoJoinWhere, page, limit int) (*DemoJoinPage, error) {
	var list []DemoInfoJoinRow
	pageObj := DemoJoinPage{Page: widget.Page{Page: page, Size: limit}}
	// 复用单表repo统一DB，无需重复绑定Model
	db := r.getGormModel(ctx)

	// Preload关联查询：表名取自shop.DemoInfoPrice.TableName
	// 回调内`demo_info_price`仅为手写SQL字符串，不参与GORM表名识别
	db = db.Joins("Price", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("`demo_info_price`.`id` as price_id, `demo_info_price`.`strategy_id`, `demo_info_price`.`price_amount`")
	})

	// 主表仅查询自身字段，严禁写入关联表字段
	db = db.Select([]string{
		"`demo_info`.`id`",
		"`demo_info`.`title`",
		"`demo_info`.`status`",
		"`demo_info`.`path`",
		"`demo_info`.`root_id`",
	})
	// 拼接连表where条件
	if w != nil {
		db = w.Build(db)
	}
	// 固定原生表名排序，禁止简写
	db = db.Order("`demo_info`.`id` DESC")
	offset := pageObj.Page.GetStart()
	// 多查一条用于判断是否还有下一页
	db = db.Offset(offset).Limit(limit + 1)

	if err := db.Find(&list).Error; err != nil {
		return pageObj.DoPage(), err
	}
	pageObj.List = list
	return pageObj.DoPage(), nil
}
